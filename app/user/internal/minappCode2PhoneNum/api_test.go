package minappCode2PhoneNum

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/siyuiot/siyu/app/user/internal/wechatAccessToken"
	"github.com/siyuiot/siyu/pkg/qpostgresql"

	_ "github.com/lib/pq"
)

func TestA(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	db, _, err := qpostgresql.InitPg("postgres://postgres:iLoveShark@192.168.0.247:32432/master?sslmode=disable&fallback_application_name=test", 10, 5)
	if err != nil {
		t.Error(err)
		return
	}
	dbRo, _, err := qpostgresql.InitPg("postgres://postgres:iLoveShark@192.168.0.247:32432/master?sslmode=disable&fallback_application_name=test", 10, 5)
	if err != nil {
		t.Error(err)
		return
	}

	wechatAccessToken.New(wechatAccessToken.Option{Log: logrus.NewEntry(log), Db: db, DbRo: dbRo})
	wat := wechatAccessToken.Instance().GetFromDbOrRemote("wxb71c87a341a6eda7")
	if wat == nil {
		t.Errorf("accessToken:%s", "")
		return
	}
	info := POST(wat.AccessToken, "089f2ee26b583f37a9bfb26b80df6e9a9fd7bf49682d7c7fddf523c38dc54eb1")
	t.Logf("info=%#+v\n", info)
}
