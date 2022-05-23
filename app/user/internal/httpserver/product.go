package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/product"
	"github.com/siyuiot/siyu/modules/qstate"
)

func (t *HttpServer) ProductList(c *gin.Context) {
	type (
		Req struct {
			Pid int `json:"pid"`
		}
		Data struct {
			Pid         int    `json:"pid"`
			Name        string `json:"name"`
			Des         string `json:"des"`
			Price       int    `json:"price"`
			PriceOrigin int    `json:"priceOrigin"`
		}
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			CustInfo  string `json:"custInfo,omitempty"` //自定义描述
			Data      []Data `json:"data"`
		}
	)
	var req Req
	var ret Ret
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	// rid := 1
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
		return
	}
	list, err := product.GetProductSkuByPid(req.Pid)
	if err != nil {
		app.Log.Error(err)
		ret.State, ret.StateInfo = State(qstate.StateFailed)
		ret.CustInfo = err.Error()
		return
	}
	for _, v := range list {
		ret.Data = append(ret.Data, Data{
			Pid:         v.Id,
			Name:        v.Name,
			Des:         v.Des,
			Price:       v.Price,
			PriceOrigin: v.PriceOrigin,
		})
	}
	ret.State, ret.StateInfo = State(qstate.StateOk)
}
