package httpserver

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/pkg/qgin"
	"github.com/siyuiot/siyu/pkg/qlog"
	"github.com/spf13/viper"
)

func (t *HttpServer) setupRouter() {
	t.gin.GET("/ping", t.Pong)

	group := t.gin.Group("/user")
	group.POST("/login", t.UserInfo)
	group.POST("/query/info", t.UserInfo)
}

func (*HttpServer) Pong(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

type HttpServer struct {
	Option
	gin *gin.Engine
}

type Option struct {
	Ctx  context.Context
	Log  *qlog.Entry
	Addr string
}

func New(o Option) error {
	t := &HttpServer{
		Option: o,
		gin:    gin.New(),
	}
	t.gin.Use(qgin.LogReq(app.Log))
	t.setupRouter()
	app.Log.Debugf("Listening and serving HTTP on %s\n", o.Addr)
	return t.gin.Run(o.Addr)
}

func State(code int) (int, string) {
	return code, viper.GetString(fmt.Sprintf("state.%d", code))
}
