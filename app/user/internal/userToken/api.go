package userToken

//用户Token

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var this *object

type Object interface {
	Upsert(info *Info) (ok bool)
	CheckAndAddTTL(token string) (int, error)
}

type Option struct {
	Log  *logrus.Entry
	Db   *sql.DB
	DbRo *sql.DB
}

type object struct {
	Option
}

func New(o Option) {
	this = &object{Option: o}
}

func Instance() *object {
	return this
}

func (o *Option) Upsert(info *Info) (ok bool) {
	_, err := o.Db.Exec(`
	INSERT INTO user_token(uid, ts, token, expire, des)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (uid)
	DO UPDATE SET ts=$2, token=$3, expire=$4, des=$5;
	`, info.Uid, info.Ts, info.Token, info.Expires, info.Des)
	if err != nil {
		o.Log.Errorf("upsert user_token err(%v) ", err)
		return false
	}
	return true
}

func (o object) queryInfo(id string) *Info {
	r := new(Info)
	var qstr string
	switch {
	case len(id) > 0:
		qstr += fmt.Sprintf(" and token = '%s'", id)
	default:
		o.Log.Error("invalid param")
		return nil
	}
	sqlstr := `
	select
	coalesce(token,'') as token,
	coalesce(ts,0) as ts,
	coalesce(uid,0) as uid,
	coalesce(expire,0) as expire,
	coalesce(des,'') as des
	from user_token
	where 1=1
	`
	sqlstr += qstr
	err := o.DbRo.QueryRow(sqlstr).Scan(&r.Token, &r.Ts, &r.Uid, &r.Expires, &r.Des)
	if err != nil {
		o.Log.Errorf("param=%s,sql=%s,err=%v", id, sqlstr, err)
		return nil
	}
	return r
}

func (o *object) CheckAndAddTTL(token string) (r int, err error) {
	now := time.Now()
	info := o.queryInfo(token)
	if info == nil {
		err := fmt.Errorf("token:%s not found", token)
		o.Log.Error(err)
		return -1, err
	}
	if info.Expires-now.Unix() < 0 {
		err := fmt.Errorf("token:%s expire at:%d,now:%d", info.Token, info.Expires, now.Unix())
		o.Log.Error(err)
		return -2, err
	}
	if info.Uid <= 0 {
		err := fmt.Errorf("token:%s err user:%d", info.Token, info.Uid)
		o.Log.Error(err)
		return -3, err
	}
	// token续命24小时
	info.Expires = info.Expires + 86400
	o.Upsert(info)
	return info.Uid, nil
}
