package httpserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/product"
	"github.com/siyuiot/siyu/app/user/internal/sim"
	"github.com/siyuiot/siyu/app/user/internal/simOrder"
	"github.com/siyuiot/siyu/app/user/internal/user"
	"github.com/siyuiot/siyu/app/user/internal/userSim"
	"github.com/siyuiot/siyu/app/user/internal/userSimLog"
	"github.com/siyuiot/siyu/app/user/internal/wechatPay"
	"github.com/siyuiot/siyu/modules/qstate"
)

func (t *HttpServer) Payment(c *gin.Context) {
	type (
		Req struct {
			AppId  string `json:"appId"`
			OpenId string `json:"openId"`
			Sid    int    `json:"sid"`
			SkuId  int    `json:"skuId"`
		}
		Ret struct {
			State     int         `json:"state"`
			StateInfo string      `json:"stateInfo"`
			CustInfo  string      `json:"custInfo,omitempty"` //自定义描述
			Data      interface{} `json:"data"`
		}
	)
	var req Req
	var ret Ret
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
		return
	}
	ruid, _ := strconv.Atoi(c.Request.Header.Get("ruid"))
	us := userSim.Instance().QueryInfo(ruid, req.Sid)
	if us == nil {
		err := fmt.Errorf("sid:%d not found", req.Sid)
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, "请先绑定sim卡"
		return
	}
	ps := product.Instance().QueryProductSkuBySkuId(req.SkuId)
	if ps == nil {
		err := fmt.Errorf("skuId:%d not found", req.SkuId)
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, "请选择套餐"
		return
	}
	now := time.Now()
	so := &simOrder.Info{
		Uid:         ruid,
		Typ:         1,
		No:          fmt.Sprintf("%d%d", 10, now.UnixNano()), // 10 sim套餐订单,
		Name:        ps.Name,
		SkuId:       ps.Id,
		Status:      simOrder.PaymentSubmit,
		AmountPrice: ps.PriceOrigin,
		DuePrice:    ps.Price,
		PayPrice:    0,
		PayChannel:  "",
		Remark:      "",
	}
	outJsApi := wechatPay.Instance().JsApi(wechatPay.InJsApi{
		AppId:          req.AppId,
		OpenId:         req.OpenId,
		OrderNo:        so.No,
		OrderPrice:     so.DuePrice,
		OrderName:      so.Name,
		OrderRemark:    so.Remark,
		OrderNotifyUrl: "https://siyu.d.blueshark.com/simOrder/payment/notify",
	})
	if outJsApi == nil {
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, "订单下单失败"
		return
	}
	r := simOrder.Instance().Insert(so)
	if r <= 0 {
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, "订单创建失败"
		return
	}
	ret.Data = outJsApi
	ret.State, ret.StateInfo = State(qstate.StateOk)
}

func (t *HttpServer) PaymentNotify(c *gin.Context) {
	type (
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
		}
	)
	var ret Ret
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	// ruid, _ := strconv.Atoi(c.Request.Header.Get("ruid"))
	transaction := wechatPay.Instance().Notify(c.Request)
	if transaction == nil {
		ret.State, ret.StateInfo = qstate.StateFailed, "支付回调Notify失败"
		return
	}
	info := simOrder.Instance().QueryInfoByNo(*transaction.OutTradeNo)
	if info == nil {
		ret.State, ret.StateInfo = qstate.StateFailed, "支付回调订单未找到"
		return
	}
	b, _ := transaction.MarshalJSON()
	r := simOrder.Instance().Update(&simOrder.Info{
		No:         *transaction.OutTradeNo,
		Status:     simOrder.PaymentSuccess,
		PayPrice:   *transaction.Amount.Total,
		PayChannel: *transaction.BankType,
		Remark:     string(b),
	})
	if len(r) <= 0 {
		ret.State, ret.StateInfo = qstate.StateFailed, "支付回调订单更新失败"
		return
	}
	go paymentSuccess(info)
	ret.State, ret.StateInfo = State(qstate.StateOk)
}

func paymentSuccess(orderInfo *simOrder.Info) {
	// 查询用户绑定sim卡信息
	us := userSim.Instance().QueryInfo(orderInfo.Uid, orderInfo.Sid)
	if us == nil {
		app.Log.Error("user not bind sim")
		return
	}
	now := time.Now()
	var orderSimDuration int = 1 // 用户购买套餐时长 单位：月=30天
	var orderSimByte int = 300   // 用户购买套餐流量 单位：G
	usNew := &userSim.Info{
		Uid:              us.Uid,
		Sid:              us.Sid,
		ServiceEndTs:     us.ServiceEndTs + int64(orderSimDuration*30*86400),
		ServiceDuration:  us.ServiceDuration + orderSimDuration,
		SimByte:          us.SimByte + orderSimByte*sim.OneG,
		SimAvailableByte: us.SimAvailableByte + orderSimByte*sim.OneG,
	}
	// 添加操作记录
	rusl := userSimLog.Instance().Insert(&userSimLog.Info{
		Ts:  now.Unix(),
		Uid: us.Uid,
		Sid: us.Sid,
		PhoneNum: func(uid int) string {
			ui := user.Instance().QueryInfo(uid, "", "", "", "")
			if ui == nil {
				return ""
			}
			return ui.PhoneNum
		}(orderInfo.Uid),
		SimNo: us.SimNo,
		Imsi:  "",
		IccId: us.Iccid,
		Remark: fmt.Sprintf(`{
			"operation": "SimOrderPaymentSuccess",
			"serviceDuration": "%d-%d",
			"serviceEndTs": "%d-%d",
			"serviceEndAt": "%s-%s",
			"simByte": "%d-%d",
			"simAvailableByte": "%d-%d"
		}`,
			us.ServiceDuration, usNew.ServiceDuration,
			us.ServiceEndTs, usNew.ServiceEndTs,
			time.Unix(us.ServiceEndTs, 0).Format("2006-01-02"), time.Unix(usNew.ServiceEndTs, 0).Format("2006-01-02"),
			us.SimByte, usNew.SimByte,
			us.SimAvailableByte, usNew.SimAvailableByte,
		),
	})
	if rusl <= 0 {
		app.Log.Error("记录写入失败")
		return
	}
	usr := userSim.Instance().Update(usNew)
	if usr <= 0 {
		app.Log.Error("用户绑定sim卡信息更新失败")
		return
	}
}
