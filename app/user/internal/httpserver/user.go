package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/modules/qstate"
)

func (t *HttpServer) UserInfo(c *gin.Context) {
	type (
		Data struct {
			Uid  int    `json:"uid"`
			Name string `json:"name"`
		}
		Req struct {
			Typ int    `json:"tpy" binding:"required"`
			Pn  string `json:"pn" binding:"required"`
		}
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			Data      Data   `json:"data"`
		}
	)
	var req Req
	var ret Ret
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		t.Log.Error(err)
		ret.State, ret.StateInfo = State(qstate.StateInvalidParameter)
		return
	}

	ret.Data = Data{}
	ret.State, ret.StateInfo = State(qstate.StateOk)
}
