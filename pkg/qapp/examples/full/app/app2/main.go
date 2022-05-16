package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/pkg/qapp"
	"github.com/siyuiot/siyu/pkg/qapp/examples/full/pkg/db"
	"github.com/siyuiot/siyu/pkg/qapp/qdebugserver"
)

const (
	appName = "app2"
	addr    = ":8088"
)

func initDB(ctx context.Context) (qapp.CleanFunc, error) {
	db.Init(ctx, appName)

	return nil, nil
}

func runHTTPServer(ctx context.Context) error {
	r := gin.Default()
	qdebugserver.RegisterGin(r, "/dev")
	return r.Run(addr)
}

func main() {
	qapp.New(appName).
		AddInitStage("initDB", initDB).
		AddDaemons(runHTTPServer).
		Run()
}
