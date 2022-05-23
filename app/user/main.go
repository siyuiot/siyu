package main

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/httpserver"
	"github.com/siyuiot/siyu/app/user/internal/product"
	"github.com/siyuiot/siyu/app/user/internal/sim"
	"github.com/siyuiot/siyu/app/user/internal/user"
	"github.com/siyuiot/siyu/app/user/internal/userSim"
	"github.com/siyuiot/siyu/app/user/internal/userSimLog"
	"github.com/siyuiot/siyu/app/user/internal/userToken"
	"github.com/siyuiot/siyu/app/user/internal/wechatAccessToken"
	"github.com/siyuiot/siyu/pkg/qapp"
	"github.com/siyuiot/siyu/pkg/qgin"
	"github.com/siyuiot/siyu/pkg/qpostgresql"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

var MASTERDB *sql.DB
var ORDERDB *gorm.DB

func preload() error {
	return nil
}

func initPgUserDb(ctx context.Context) (qapp.CleanFunc, error) {
	dbUrl := viper.GetString("pg.user.url")
	if qgin.IsDev() {
		dbUrl = qgin.DbUserUrl
	}
	userdb, _, err := qpostgresql.InitPg(dbUrl, viper.GetInt("pg.user.maxopen"), viper.GetInt("pg.user.maxidle"))
	if err != nil {
		app.Log.Errorf("Fatal initPg: %s \n", err)
		return nil, err
	}
	MASTERDB = userdb
	return nil, nil
}

func initPgOrderDb(ctx context.Context) (qapp.CleanFunc, error) {
	dbUrl := viper.GetString("pg.user.url")
	if qgin.IsDev() {
		dbUrl = qgin.DbUserUrl
	}
	orderDb, err := qpostgresql.InitGorm(dbUrl, viper.GetInt("pg.user.maxopen"), viper.GetInt("pg.user.maxidle"))
	if err != nil {
		app.Log.Errorf("Fatal initPg: %s \n", err)
		return nil, err
	}
	ORDERDB = orderDb
	return nil, nil
}

func initModules(ctx context.Context) (qapp.CleanFunc, error) {
	user.New(user.Option{Log: app.Log, Db: MASTERDB, DbRo: MASTERDB})
	userToken.New(userToken.Option{Log: app.Log, Db: MASTERDB, DbRo: MASTERDB})
	wechatAccessToken.New(wechatAccessToken.Option{Log: app.Log, Db: MASTERDB, DbRo: MASTERDB})
	product.New(product.Option{Log: app.Log, DBORM: ORDERDB})
	sim.New(sim.Option{Log: app.Log, Db: MASTERDB, DbRo: MASTERDB})
	userSim.New(userSim.Option{Log: app.Log, Db: MASTERDB, DbRo: MASTERDB})
	userSimLog.New(userSimLog.Option{Log: app.Log, Db: MASTERDB, DbRo: MASTERDB})
	return func(ctx context.Context) {}, nil
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
		AddInitStage("initPgUserDb", initPgUserDb).
		AddInitStage("initPgOrderDb", initPgOrderDb).
		AddInitStage("modules", initModules).
		AddDaemons(runServer).
		Run()
}
