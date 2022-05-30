package wechatPay

// https://github.com/wechatpay-apiv3/wechatpay-go

// 微信支付

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/siyuiot/siyu/pkg/qlog"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var this *object

// type Object interface {
// 	JsApi()
// 	Notify()
// }

type Option struct {
	Log                        *qlog.Entry
	Db                         *sql.DB
	DbRo                       *sql.DB
	MchID                      string // 微信支付商户号
	mchCertificateSerialNumber string
	mchAPIV3Key                string
	mchPrivateKey              string
	mchCert                    string
	client                     *core.Client
}

type object struct {
	Option
}

func New(o Option) {
	this = &object{
		Option: o,
	}
}

func Instance() *object {
	return this
}

func (o object) queryWechatPayMch() (mchCertificateSerialNumber, mchAPIv3Key, mchPrivateKey, mchCert string) {
	var qstr string
	switch {
	case len(o.MchID) > 0:
		qstr += fmt.Sprintf(" and mch_id = '%s'", o.MchID)
	default:
		o.Log.Error("invalid param")
		return
	}
	sqlstr := `
	select
	coalesce(mch_certificate_serial_number,'') as mch_certificate_serial_number,
	coalesce(mch_API_v3_key,'') as mch_API_v3_key,
	coalesce(mch_private_key,'') as mch_private_key,
	coalesce(mch_cert,'') as mch_cert
	from wechat_pay_mch
	where 1=1
	`
	sqlstr += qstr
	// o.Log.Debug(sqlstr)
	err := o.DbRo.QueryRow(sqlstr).Scan(&mchCertificateSerialNumber, &mchAPIv3Key, &mchPrivateKey, &mchCert)
	if err != nil {
		o.Log.Errorf("param=%s,sql=%s,err=%v", o.MchID, sqlstr, err)
		return
	}
	return
}

func (o object) NewWechatPayClient() *core.Client {
	o.mchCertificateSerialNumber, o.mchAPIV3Key, o.mchPrivateKey, o.mchCert = o.queryWechatPayMch()
	// 使用 utils 提供的函数加载商户私钥，商户私钥会用来生成请求的签名
	privateKey, err := utils.LoadPrivateKey(o.mchPrivateKey)
	if err != nil {
		o.Log.Errorf("load merchant private key error:%s", err)
		return nil
	}
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(o.MchID, o.mchCertificateSerialNumber, privateKey, o.mchAPIV3Key),
	}
	client, err := core.NewClient(context.Background(), opts...)
	if err != nil {
		o.Log.Errorf("new wechat pay client err:%s", err)
		return nil
	}
	return client
}

func (o object) JsApi(in InJsApi) *jsapi.PrepayWithRequestPaymentResponse {
	if o.client == nil {
		o.client = o.NewWechatPayClient()
	}
	// JSAPI下单
	svc := jsapi.JsapiApiService{Client: o.client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, _, err := svc.PrepayWithRequestPayment(context.Background(), jsapi.PrepayRequest{
		Appid:       core.String(in.AppId),
		Mchid:       core.String(o.MchID),
		Description: core.String(in.OrderName),
		OutTradeNo:  core.String(in.OrderNo),
		Attach:      core.String(in.OrderRemark),
		NotifyUrl:   core.String(in.OrderNotifyUrl),
		Amount: &jsapi.Amount{
			Total: core.Int64(in.OrderPrice),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(in.OpenId),
		},
	})
	if err != nil {
		o.Log.Error(err)
		return nil
	}
	return resp
}

// https://github.com/wechatpay-apiv3/wechatpay-go/blob/main/FAQ.md#%E4%B8%BA%E4%BB%80%E4%B9%88%E6%94%B6%E5%88%B0%E5%BA%94%E7%AD%94%E4%B8%AD%E7%9A%84%E8%AF%81%E4%B9%A6%E5%BA%8F%E5%88%97%E5%8F%B7%E5%92%8C%E5%8F%91%E8%B5%B7%E8%AF%B7%E6%B1%82%E7%9A%84%E8%AF%81%E4%B9%A6%E5%BA%8F%E5%88%97%E5%8F%B7%E4%B8%8D%E4%B8%80%E8%87%B4
func (o object) certDownload() {
	o.mchCertificateSerialNumber, o.mchAPIV3Key, o.mchPrivateKey, o.mchCert = o.queryWechatPayMch()
	// 使用 utils 提供的函数加载证书
	privateKey, err := utils.LoadPrivateKey(o.mchPrivateKey)
	if err != nil {
		o.Log.Error(err)
		return
	}
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, privateKey, o.mchCertificateSerialNumber, o.MchID, o.mchAPIV3Key)
	if err != nil {
		o.Log.Error(err)
		return
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	m := downloader.MgrInstance().GetCertificateMap(ctx, o.MchID)
	// cert, b := certificateVisitor.Get(context.Background(), o.mchCertificateSerialNumber)
	o.Log.Debugf("cert:%#+v,b:%v", m, err)
}

// request.Body在ParseNotifyRequest前不能读取
func (o object) Notify(request *http.Request) *payments.Transaction { //request *http.Request
	o.mchCertificateSerialNumber, o.mchAPIV3Key, o.mchPrivateKey, o.mchCert = o.queryWechatPayMch()
	// 使用 utils 提供的函数加载证书
	privateKey, err := utils.LoadPrivateKey(o.mchPrivateKey)
	if err != nil {
		o.Log.Error(err)
		return nil
	}
	ctx := context.Background()
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(ctx, privateKey, o.mchCertificateSerialNumber, o.MchID, o.mchAPIV3Key)
	if err != nil {
		o.Log.Error(err)
		return nil
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(o.MchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(o.mchAPIV3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	// 4. 验签与解密
	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), request, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		o.Log.Error(err)
		return nil
	}
	// 处理通知内容
	o.Log.Debugf("url=/simOrder/payment/notify req=%#+v", notifyReq)
	if notifyReq.EventType == "TRANSACTION.SUCCESS" {
		return transaction
	}
	return nil
}
