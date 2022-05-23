package user

import (
	"fmt"
	"testing"
	"time"

	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/pkg/qpostgresql"
)

func initInstance() Object {
	db, _, err := qpostgresql.InitPg("postgres://postgres:iLoveShark@192.168.0.247:32432/bsmaster?sslmode=disable&fallback_application_name=devicedatasvc", 10, 5)
	if err != nil {
		app.Log.Error(err)
		return nil
	}
	return New(Option{
		Log: app.Log,
		Db:  db,
	})
}

func TestInsert(t *testing.T) {
	now := time.Now()
	entry := initInstance()
	uid := entry.Insert(&UserInfo{
		UserId:      0,
		PhoneNum:    "",
		PhoneArea:   "86",
		Account:     "",
		Email:       "huangzhifeng@sharkgulf.com",
		Password:    "111111",
		PwdSalt:     "18616854987",
		NickName:    "tobeno.1",
		RealName:    "黄志峰",
		Gender:      1,
		Birthday:    "1982-05-04",
		IdNo:        "4211811",
		Icon:        "aaa",
		Location:    "",
		PerSign:     "",
		Completion:  100,
		Home:        "{}",
		RegType:     1,
		Tz:          "",
		App:         "blueshark",
		CreatedTime: now,
		UpdatedTime: now,
		// GeneralSetup: "",
	})
	fmt.Printf("%+v\n", uid)
}

// func TestInsert(t *testing.T) {
// 	entry := initInstance()
// 	uid := entry.Insert(&Info{
// 		UserId:          999,
// 		Partner:         "qq",
// 		PartnerOpenid:   "aaa",
// 		PartnerUid:      "aaa",
// 		PartnerNickName: "aaa",
// 		BindedTime:      time.Now().Unix(),
// 		App:             user.Blueshark,
// 	})
// 	fmt.Printf("%+v\n", uid)
// }

// func TestQuery(t *testing.T) {
// 	entry := initInstance()
// 	info := entry.QueryInfoByUids(999, "qq", user.Blueshark)
// 	fmt.Printf("%+v\n", info)
// 	info = entry.QueryInfoPartnerOpenid("qq", "2BA47C96B957E8D28D08B91B1DD0BD87", user.Blueshark)
// 	fmt.Printf("%+v\n", info)
// }

func TestQueryInfo(t *testing.T) {
	info := initInstance().QueryInfo(180, "", "", "", "")
	if info == nil {
		fmt.Println("error")
	} else {
		fmt.Printf("%+v", *info)
	}
}
