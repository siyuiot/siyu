package wechatPay

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/siyuiot/siyu/pkg/qpostgresql"

	_ "github.com/lib/pq"
)

func initInstance() {
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
		return
	}
	dbRo, _, err := qpostgresql.InitPg("postgres://postgres:iLoveShark@192.168.0.247:32432/master?sslmode=disable&fallback_application_name=test", 10, 5)
	if err != nil {
		log.Error(err)
		return
	}
	New(Option{
		Log:    logrus.NewEntry(log),
		Db:     db,
		DbRo:   dbRo,
		MchID:  "1626436719",
		client: nil,
	})
}

func TestJsApi(t *testing.T) {
	initInstance()
	Instance().JsApi(InJsApi{
		AppId:          "wxb71c87a341a6eda7",
		OpenId:         "",
		OrderNo:        "10-20220528221400",
		OrderPrice:     100,
		OrderName:      "测试订单名称",
		OrderRemark:    "测试订单备注",
		OrderNotifyUrl: "https://siyu.d.blueshark.com/simOrder/payment/notify",
	})
}

func TestCD(t *testing.T) {
	initInstance()
	Instance().certDownload()
}
