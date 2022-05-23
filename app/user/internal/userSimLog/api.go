package userSimLog

// 用户绑定sim卡log

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/siyuiot/siyu/pkg/qlog"
)

var this *object

type Object interface {
	Insert(*Info) int
}

type Option struct {
	Log  *qlog.Entry
	Db   *sql.DB
	DbRo *sql.DB
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

func (o object) Insert(i *Info) (r int) {
	now := time.Now()
	i.Created = now.Unix()
	sqlstr := `insert into user_sim_log(ts,uid,sid,phone_num,sim_no,imsi,iccid,remark,created)
	values(%d,%d,%d,'%s','%s','%s','%s','%s',%d)
	returning ts`
	sqlstr = fmt.Sprintf(sqlstr, i.Ts, i.Uid, i.Sid, i.PhoneNum, i.SimNo, i.Imsi, i.IccId, i.Remark, i.Created)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", i, sqlstr, err)
		return
	}
	return
}
