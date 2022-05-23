package product

// 产品：用户服务

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/siyuiot/siyu/pkg/qhttp"
	"github.com/siyuiot/siyu/pkg/qlog"
)

var this *object

var orderDB *gorm.DB

type Object interface {
	// Update(*Info) int
	// QueryInfoFromDb(appId string) *Info
	QueryInfoFromDbOrRemote(appId string) *Info
}

type Option struct {
	Log   *qlog.Entry
	Db    *sql.DB
	DbRo  *sql.DB
	DBORM *gorm.DB
}

type object struct {
	Option
}

func New(o Option) {
	this = &object{
		Option: o,
	}
}

func Instance() *object {
	return this
}

func (o object) update(i *Info) (r string) {
	now := time.Now().UTC()
	sqlstr := fmt.Sprintf("update wechat_access_token set updated = %d,", now.Unix())
	if len(i.AppId) <= 0 {
		o.Log.Error("invalid param")
		return
	}
	if len(i.AccessToken) > 0 {
		sqlstr += fmt.Sprintf("access_token = '%s',", i.AccessToken)
	}
	if i.ExpiresIn > 0 {
		sqlstr += fmt.Sprintf("expires_in = '%d',", i.ExpiresIn)
	}
	if i.ExpiresAt > 0 {
		sqlstr += fmt.Sprintf("expires_at = '%d',", i.ExpiresAt)
	}
	sqlstr = strings.TrimRight(sqlstr, ",")
	sqlstr += fmt.Sprintf("where app_id = '%s' returning app_id", i.AppId)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", i, sqlstr, err)
		return
	}
	return
}

func (o object) queryInfo(appId string) *Info {
	r := new(Info)
	var qstr string
	switch {
	case len(appId) > 0:
		qstr += fmt.Sprintf(" and app_id = '%s'", appId)
	default:
		o.Log.Error("invalid param")
		return nil
	}
	sqlstr := `
	select
	coalesce(app_id,'') as app_id,
	coalesce(secret,'') as secret,
	coalesce(access_token,'') as access_token,
	coalesce(expires_in,0) as expires_in,
	coalesce(expires_at,0) as expires_at,
	coalesce(created,0) as created,
	coalesce(updated,0) as updated
	from wechat_access_token
	where 1=1
	`
	sqlstr += qstr
	err := o.DbRo.QueryRow(sqlstr).Scan(&r.AppId, &r.Secret, &r.AccessToken, &r.ExpiresIn, &r.ExpiresAt, &r.Created, &r.Updated)
	if err != nil {
		o.Log.Errorf("param=%s,sql=%s,err=%v", appId, sqlstr, err)
		return nil
	}
	return r
}

func (o object) queryInfoFromDb(appId string) *Info {
	return o.queryInfo(appId)
}

func (o object) GetFromDbOrRemote(appId string) *Info {
	now := time.Now().UTC()
	info := o.queryInfoFromDb(appId)
	if info == nil {
		o.Log.Error("appId is not config")
		return nil
	}
	// accessToken有值
	// accessToken没过期
	if len(info.AccessToken) > 0 && info.ExpiresAt > now.Unix() {
		o.Log.Info("accessToken from db")
		return info
	}

	// 从微信获取accessToken
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", info.AppId, info.Secret)
	type WechatAccessToken struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	var wat WechatAccessToken
	resp, err := qhttp.GetJSON(url, &wat)
	if err != nil {
		o.Log.Error(err)
		return nil
	}
	if len(wat.AccessToken) <= 0 {
		o.Log.Errorf("wechatAccessToken:%s,err:%s", wat.AccessToken, resp)
		return nil
	}

	// 更新accessToken到DB
	rid := o.update(&Info{
		AppId:       appId,
		AccessToken: wat.AccessToken,
		ExpiresIn:   wat.ExpiresIn,
		ExpiresAt:   now.Unix() + int64(wat.ExpiresIn),
	})
	if len(rid) <= 0 {
		o.Log.Errorf("wechatAccessToken update error")
		return nil
	}
	o.Log.Info("accessToken from remote")
	return o.queryInfoFromDb(rid)
}
