package main

import (
	"context"
	"database/sql"

	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/httpserver"
	"github.com/siyuiot/siyu/pkg/qapp"
	"github.com/siyuiot/siyu/pkg/qgin"
	"github.com/siyuiot/siyu/pkg/qpostgresql"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

var USERDB *sql.DB

func preload() error {
	return nil
}

func initPgDb(ctx context.Context) (qapp.CleanFunc, error) {
	dbUrl := viper.GetString("pg.user.url")
	if qgin.IsDev() {
		dbUrl = qgin.DbUserUrl
	}
	userdb, _, err := qpostgresql.InitPg(dbUrl, viper.GetInt("pg.user.maxopen"), viper.GetInt("pg.user.maxidle"))
	if err != nil {
		app.Log.Errorf("Fatal initPg: %s \n", err)
		return nil, err
	}
	USERDB = userdb
	return nil, nil
}

func runServer(ctx context.Context) error {
	httpserver.New(httpserver.Option{
		Ctx:  ctx,
		Log:  app.Log,
		Addr: viper.GetString("app.addr"),
	})
	return nil
}

func main() {
	qapp.New(app.Name, qapp.WithPreload(preload)).
		AddInitStage("database", initPgDb).
		AddDaemons(runServer).
		Run()
}
