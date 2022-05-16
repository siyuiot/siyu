package qgin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/modules/qstate"
	"github.com/spf13/viper"
)

type HttpContext interface {
	ShouldBindJSON(obj interface{}) error
	AbortWithStatusJSON(code int, jsonObj interface{})
}

func httpCtx(c *gin.Context) HttpContext {
	return &httpContext{c: c}
}

func checkJsonParam(c HttpContext, params interface{}) bool {
	err := c.ShouldBindJSON(params)
	if err != nil {
		info := viper.GetString(fmt.Sprint("state.", qstate.StateInvalidParameter))
		if info == "" {
			info = "无效参数"
		}
		ret := qstate.CommonResponse{
			State:     qstate.StateInvalidParameter,
			StateInfo: info,
			CustInfo:  err.Error(),
			Seq:       0,
		}
		c.AbortWithStatusJSON(http.StatusOK, ret)
		return false
	}
	return true
}

func sendOk(c HttpContext, data interface{}) {
	send(c, qstate.StateOk, data, "", "", 0)
}

func sendOkAndSeq(c HttpContext, data interface{}, seq int) {
	send(c, qstate.StateOk, data, "", "", seq)
}

// state func
func state(s qstate.State) string {
	info := viper.GetString(fmt.Sprintf("state.", s))
	if info == "" {
		info = qstate.StateStr(s)
	}
	return info
}

func send(c HttpContext, code qstate.State, data interface{}, stateInfo, CustomInfo string, seq int) {
	info := stateInfo
	if info == "" {
		info = state(code)
	}
	ret := qstate.CommonResponse{
		State:     code,
		StateInfo: info,
		CustInfo:  CustomInfo,
		Seq:       seq,
		Data:      data,
	}
	c.AbortWithStatusJSON(http.StatusOK, ret)
}

func sendFailed(c HttpContext, data interface{}) {
	send(c, qstate.StateFailed, data, "", "", 0)
}

func sendPureFailed(c HttpContext) {
	send(c, qstate.StateFailed, nil, "", "", 0)
}

func sendNotFound(c HttpContext) {
	send(c, qstate.StateNoRecord, nil, "", "", 0)
}
func sendDuplicate(c HttpContext) {
	send(c, qstate.StateDuplicate, nil, "重复创建", "", 0)
}
func sendFailedWithErr(c HttpContext, err error) {
	send(c, qstate.StateFailed, nil, err.Error(), "", 0)
}
func sendFailedWithMsg(c HttpContext, errorMsg string) {
	send(c, qstate.StateFailed, nil, errorMsg, "", 0)
}
func sendPureOk(c HttpContext) {
	sendOk(c, "")
}

func sendParamError(c HttpContext, msg string) {
	send(c, qstate.StateInvalidParameter, nil, msg, "", 0)
}

func notLogin(c HttpContext) {
	send(c, qstate.StateNotAuth, nil, "", "", 0)
}

func sendList(c HttpContext, list interface{}, total int) {
	sendOk(c, map[string]interface{}{
		"total": total,
		"list":  list,
	})
}
