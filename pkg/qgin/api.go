package qgin

import (
	"github.com/siyuiot/siyu/modules/qstate"
	"github.com/siyuiot/siyu/pkg/qerrors"
	"github.com/siyuiot/siyu/pkg/qhttp"

	"github.com/gin-gonic/gin"
)

type httpContext struct {
	c *gin.Context
}

func (h *httpContext) ShouldBindJSON(obj interface{}) error {
	return h.c.ShouldBindJSON(obj)
}

func (h *httpContext) AbortWithStatusJSON(code int, jsonObj interface{}) {
	h.c.AbortWithStatusJSON(code, jsonObj)
}

func CheckJsonParam(c *gin.Context, params interface{}) bool {
	return checkJsonParam(httpCtx(c), params)
}

func CheckJsonParamErr(c *gin.Context, params interface{}) error {
	err := c.ShouldBindJSON(params)
	if err != nil {
		return qerrors.New(qstate.StateInvalidParameter, "参数错误")
	}
	return nil
}

func SendOk(c *gin.Context, data interface{}) {
	sendOk(httpCtx(c), data)
}

func SendOkAndSeq(c *gin.Context, data interface{}, seq int) {
	sendOkAndSeq(httpCtx(c), data, seq)
}

func Send(c *gin.Context, code qstate.State, data interface{}, stateInfo, customInfo string) {
	send(httpCtx(c), code, data, stateInfo, customInfo, 0)
}

func SendFailed(c *gin.Context, data interface{}) {
	sendFailed(httpCtx(c), data)
}

func SendCode(c *gin.Context, code qstate.State) {
	send(httpCtx(c), code, nil, "", "", 0)
}

func SendCodeAndData(c *gin.Context, code qstate.State, data interface{}) {
	send(httpCtx(c), code, data, "", "", 0)
}
func SendCodeAndDataSeq(c *gin.Context, code qstate.State, data interface{}, seq int) {
	send(httpCtx(c), code, data, "", "", seq)
}

func SendPureFailed(c *gin.Context) {
	sendPureFailed(httpCtx(c))
}

func SendNotFound(c *gin.Context) {
	sendNotFound(httpCtx(c))
}
func SendDuplicate(c *gin.Context) {
	sendDuplicate(httpCtx(c))
}

func SendFailedWithErr(c *gin.Context, err error) {
	sendFailedWithErr(httpCtx(c), err)
}

func SendCodeWithMsg(c *gin.Context, code qstate.State, errorMsg string) {
	send(httpCtx(c), code, "", errorMsg, "", 0)
}

func SendCodeWithErr(c *gin.Context, code qstate.State, err error) {
	send(httpCtx(c), code, "", err.Error(), "", 0)
}

func SendFailedWithMsg(c *gin.Context, errorMsg string) {
	sendFailedWithMsg(httpCtx(c), errorMsg)
}
func SendPureOk(c *gin.Context) {
	sendPureOk(httpCtx(c))
}

func SendParamError(c *gin.Context) {
	sendParamError(httpCtx(c), "")
}
func SendParamErrorWithMsg(c *gin.Context, msg string) {
	sendParamError(httpCtx(c), msg)
}

func NotLogin(c *gin.Context) {
	notLogin(httpCtx(c))
}

func SendList(c *gin.Context, list interface{}, total int) {
	sendList(httpCtx(c), list, total)
}

func SendFile(c *gin.Context, data []byte, filename string) {
	qhttp.Download(data, c.Writer, c.Request, filename)
}
