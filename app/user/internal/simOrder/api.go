package simOrder

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/pkg/qlog"
)

var this *object

type Object interface {
	Insert(*Info) int
	Update(i *Info) (r string)
	QueryInfoByNo(no string) *Info
	QueryList(uid int) (int, []Info)
	Delete(uid, sid int) *Info
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
	i.Updated = now.Unix()
	sqlstr := `insert into sim_order(uid,sid,name,no,typ,sku_id,status,amount_price,due_price,pay_channel,remark,created,updated)
	values(%d,%d,'%s','%s',%d,%d,%d,%d,%d,'%s','%s',%d,%d)
	returning oid;`
	sqlstr = fmt.Sprintf(sqlstr, i.Uid, i.Sid, i.Name, i.No, i.Typ, i.SkuId, i.Status, i.AmountPrice, i.DuePrice, i.PayChannel, i.Remark, i.Created, i.Updated)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", i, sqlstr, err)
		return
	}
	return
}

func (o object) Update(i *Info) (r string) {
	now := time.Now()
	sqlstr := fmt.Sprintf("update sim_order set updated = %d,", now.Unix())
	if len(i.No) <= 0 {
		o.Log.Error("invalid param")
		return
	}
	if i.Status > 0 {
		sqlstr += fmt.Sprintf("status = %d,", i.Status)
	}
	if i.PayPrice > 0 {
		sqlstr += fmt.Sprintf("pay_price = %d,", i.PayPrice)
	}
	if len(i.PayChannel) > 0 {
		sqlstr += fmt.Sprintf("pay_channel = '%s',", i.PayChannel)
	}
	if len(i.Remark) > 0 {
		sqlstr += fmt.Sprintf("remark = '%s',", i.Remark)
	}
	sqlstr = strings.TrimRight(sqlstr, ",")
	sqlstr += fmt.Sprintf("where no = '%s' returning no", i.No)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", i, sqlstr, err)
		return
	}
	return
}

func (o object) QueryList(uid int) (total int, list []Info) {
	list = make([]Info, 0)
	var qstr string
	switch {
	case uid > 0:
		qstr += fmt.Sprintf(` and uid = %d `, uid)
	default:
		o.Log.Error("invalid param")
		return
	}
	sqlcount := `
	select count(1) from user_sim t
	where 1=1
	`
	sqlcount += qstr
	// o.Log.Debug(sqlcount)
	err := o.Db.QueryRow(sqlcount).Scan(&total)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", uid, sqlcount, err)
		return
	}
	sqlstr := `
	select
	coalesce(uid,0) as uid,
	coalesce(sid,0) as sid,
	coalesce(sim_no,'') as sim_no,
	coalesce(sim_byte,0) as sim_byte,
	coalesce(sim_available_byte,0) as sim_available_byte,
	coalesce(bind_ts,0) as bind_ts,
	coalesce(service_end_ts,0) as service_end_ts,
	coalesce(remark,'') as remark,
	coalesce(created,0) as created,
	coalesce(updated,0) as updated
	from user_sim
	where 1=1
	`
	sqlstr += qstr
	// o.Log.Debug(sqlstr)
	rows, err := o.Db.Query(sqlstr)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", uid, sqlstr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var r = Info{}
		err = rows.Scan(&r.Uid, &r.Remark, &r.Created, &r.Updated)
		if err != nil {
			o.Log.Errorf("param=%+v,sql=%s,err=%v", uid, sqlstr, err)
			continue
		}
		list = append(list, r)
	}
	return
}

func (o object) queryInfo(no string, uid, sid int) *Info {
	r := new(Info)
	var qstr string
	if uid > 0 {
		qstr += fmt.Sprintf(" and uid = %d ", uid)
	}
	if sid > 0 {
		qstr += fmt.Sprintf(" and sid = %d ", sid)
	}
	if len(no) > 0 {
		qstr += fmt.Sprintf(" and no = '%s' ", no)
	}
	sqlstr := `
	select
	coalesce(oid,0) as oid,
	coalesce(uid,0) as uid,
	coalesce(sid,0) as sid,
	coalesce(no,'') as no,
	coalesce(name,'') as name,
	coalesce(typ,0) as typ,
	coalesce(sku_id,0) as sku_id,
	coalesce(status,0) as status,
	coalesce(remark,'') as remark,
	coalesce(created,0) as created,
	coalesce(updated,0) as updated
	from sim_order
	where 1=1
	`
	sqlstr += qstr
	app.Log.Debug(sqlstr)
	err := o.DbRo.QueryRow(sqlstr).Scan(&r.Oid, &r.Uid, &r.Sid, &r.No, &r.Name, &r.Typ, &r.SkuId, &r.Status, &r.Remark, &r.Created, &r.Updated)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", uid, sqlstr, err)
		return nil
	}
	return r
}

func (o object) QueryInfoByNo(no string) *Info {
	return o.queryInfo(no, 0, 0)
}

func (o object) Delete(uid, sid int) (r int) {
	sqlstr := `delete from user_sim where uid = $1 and sid = $2`
	res, err := o.Db.Exec(sqlstr, uid, sid)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", uid, sqlstr, err)
		return
	}
	ra, _ := res.RowsAffected()
	if ra != 1 {
		o.Log.Warnf("param=%d,sql=%s,rowsAftected=%d", uid, sqlstr, ra)
		return int(ra)
	}
	return
}
