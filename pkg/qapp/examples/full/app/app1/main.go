package main

import (
	"context"

	"github.com/siyuiot/siyu/pkg/qapp"
	"github.com/siyuiot/siyu/pkg/qapp/examples/full/pkg/db"
	"github.com/siyuiot/siyu/pkg/qapp/examples/full/pkg/httpsrv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	appName = "app1"
	addr    = ":8080"
)

func preload() error {
	pflag.String("token", "app1", "appname")
	pflag.String("addr", ":8080", "listen address")

	//viper.Set("file", "app.yml")
	return nil
}

func initDB(ctx context.Context) (qapp.CleanFunc, error) {
	db.Init(ctx, appName)

	return nil, nil
}

func runHTTPServer(ctx context.Context) error {
	return httpsrv.Run(ctx, viper.GetString("token"), viper.GetString("addr"))
}

func onConfigChange() {
	httpsrv.Restart(viper.GetString("token"), viper.GetString("addr"))
}

func main() {
	qapp.New(appName, qapp.WithPreload(preload), qapp.WithConfigChanged(onConfigChange)).
		AddInitStage("initDB", initDB).
		AddDaemons(runHTTPServer).
		Run()
}
