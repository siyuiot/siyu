package wechatPay

import "github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"

type InJsApi struct {
	AppId          string
	OpenId         string
	OrderNo        string
	OrderPrice     int64
	OrderName      string
	OrderRemark    string
	OrderNotifyUrl string
}

type OutJsApi jsapi.PrepayWithRequestPaymentResponse
