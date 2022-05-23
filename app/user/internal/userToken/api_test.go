package userToken

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/siyuiot/siyu/pkg/qpostgresql"

	_ "github.com/lib/pq"
)

func initInstance() Object {
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

func TestUpserToken(t *testing.T) {
	now := time.Now()
	ok := initInstance().Upsert(&Info{
		Uid:     1,
		Token:   "aa",
		Ts:      now.Unix(),
		Des:     "a",
		Expires: now.Add(time.Hour * 24).Unix(),
	})
	if ok {
		fmt.Println("ok")
	} else {
		t.Fail()
		fmt.Println("false")
	}
}
