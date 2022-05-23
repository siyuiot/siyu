package userSimLog

import (
	"fmt"
	"testing"
	"time"

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

func TestInsert(t *testing.T) {
	now := time.Now()
	entry := initInstance()
	info := Info{
		Ts:       now.Unix(),
		Uid:      1,
		Sid:      1,
		PhoneNum: "18616854987",
		SimNo:    "18616854987",
		Imsi:     "18616854987",
		IccId:    "18616854987",
		Remark:   "测试",
	}
	r := entry.Insert(&info)
	fmt.Printf("info=%#+v\n", r)
}
