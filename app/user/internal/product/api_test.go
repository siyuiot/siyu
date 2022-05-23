package product

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/siyuiot/siyu/pkg/qpostgresql"

	_ "github.com/lib/pq"
)

func initInstance() *object {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	db, _, err := qpostgresql.InitPg("postgres://postgres:iLoveShark@192.168.0.247:32432/master?sslmode=disable&fallback_application_name=test", 10, 5)
	if err != nil {
		log.Error(err)
		return nil
	}
	dbRo, _, err := qpostgresql.InitPg("postgres://postgres:iLoveShark@192.168.0.247:32432/master?sslmode=disable&fallback_application_name=test", 10, 5)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &object{Option{
		Log:  logrus.NewEntry(log),
		Db:   db,
		DbRo: dbRo,
	}}
}

func TestQueryInfo(t *testing.T) {
	entry := initInstance()
	info := entry.GetFromDbOrRemote("wxb71c87a341a6eda7")
	fmt.Printf("info=%#+v\n", info)
}
