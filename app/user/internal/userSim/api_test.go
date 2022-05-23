package userSim

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
		Uid:              1,
		Sid:              2,
		SimNo:            "12345678901",
		SimByte:          307200,
		SimAvailableByte: 307200,
		BindTs:           now.Unix(),
		ServiceEndTs:     now.Unix() + 86400,
		Remark:           "测试",
	}
	r := entry.Insert(&info)
	fmt.Printf("info=%#+v,now=%v\n", r, now)
}

func TestQueryInfo(t *testing.T) {
	entry := initInstance()
	info := entry.QueryInfo(1, 1)
	fmt.Printf("info=%#+v,bindTime=%s\n", info, time.Unix(info.BindTs, 0))
}

func TestDelete(t *testing.T) {
	entry := initInstance()
	info := entry.Delete(1, 1)
	fmt.Printf("info=%#+v\n", info)
}
