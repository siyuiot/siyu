package httpserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/userToken"
	"github.com/siyuiot/siyu/modules/qstate"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m = map[string]int{
			"/user/login/minapp":       1,
			"/simOrder/payment/notify": 1,
		}
		type Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			CustInfo  string `json:"custInfo,omitempty"` //自定义描述
		}
		var ret Ret
		c.Request.Header.Del("ruid") // 删除 ruid,防止前台注入
		ruid := -999
		defer func() {
			c.Request.Header.Add("ruid", strconv.Itoa(ruid))
			if ret.State == 0 {
				c.Next()
			} else {
				c.JSON(http.StatusOK, ret)
				c.Abort()
			}
		}()
		headToken := c.Request.Header.Get("Authorization") // {"Authorization":"Bearer ${token}"}
		if _, ok := m[c.Request.URL.String()]; ok {
			app.Log.Infof("path:%s not need auth token", c.Request.URL.String())
			return
		}
		if len(headToken) < 7 {
			err := fmt.Errorf("header token:%s err", headToken)
			ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
			return
		}
		clientToken := headToken[7:]
		r, err := userToken.Instance().CheckAndAddTTL(clientToken)
		if err != nil {
			app.Log.Error(err)
			ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
			return
		}
		ruid = r
	}
}
