package qgin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/pkg/qerrors"
	"github.com/siyuiot/siyu/pkg/qlog"
)

// 统一包装错误处理
func Wrap(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		reply, err := h(c)
		fromError := qerrors.FromError(err)
		if fromError == nil {
			sendOk(c, reply)
		} else {
			qlog.WithField(MiddleReqId, GetRequestId(c)).WithField(MiddleUidKey, getuid(c)).Errorf("%+v", fromError)
			var metadata = ""
			if len(fromError.Metadata) > 0 {
				metadata = fmt.Sprintf("%v", fromError.Metadata)
			}
			Send(c, fromError.State, nil, fromError.StateInfo, metadata)
		}
	}
}
