package httpserver

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/user"
	"github.com/siyuiot/siyu/pkg/qgin"
	"github.com/siyuiot/siyu/pkg/qlog"
	"github.com/spf13/viper"
)

func (t *HttpServer) setupRouter() {
	t.gin.GET("/ping", t.Pong)

	group := t.gin.Group("/user")
	group.POST("/login/minapp", t.LoginMinapp)
	group.POST("/get/info", t.Pong)

	product := t.gin.Group("/product")
	product.POST("/get/list", t.ProductList)
	product.POST("/get/info", t.Pong)

	// 用户sim卡
	service := t.gin.Group("/userSim")
	service.POST("/bind", t.UserSimBind)     // 绑定sim
	service.POST("/get/list", t.UserSimList) // 用户sim信息
	service.POST("/get/info", t.UserSimInfo) // 用户sim信息

}

func (*HttpServer) Pong(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

type HttpServer struct {
	Option
	gin *gin.Engine
}

type Option struct {
	Ctx       context.Context
	Log       *qlog.Entry
	Addr      string
	UserEntry user.Object
}

func New(o Option) error {
	t := &HttpServer{
		Option: o,
		gin:    gin.New(),
	}
	t.gin.Use(qgin.LogReq(app.Log))
	t.gin.Use(Auth())
	t.setupRouter()
	app.Log.Debugf("Listening and serving HTTP on %s\n", o.Addr)
	return t.gin.Run(o.Addr)
}

func State(code int) (int, string) {
	return code, viper.GetString(fmt.Sprintf("state.%d", code))
}
