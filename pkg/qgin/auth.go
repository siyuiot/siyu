package qgin

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// TODO 鉴权
		uid := GetUid(c)
		if uid <= 0 {
			err := errors.New("请先登录")
			logrus.Debugf("url=%s uid=%d requestId=%s", c.Request.URL.String(), uid, GetRequestId(c))
			SendFailedWithErr(c, err)
			return
		}
		c.Next()
	}
}
