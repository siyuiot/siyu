package userSim

// 用户绑定sim卡

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/siyuiot/siyu/pkg/qlog"
)

var this *object

type Object interface {
	Insert(*Info) int
	Update(*Info) int
	QueryInfo(uid, sid int) *Info
	QueryInfoBySid(sid int) *Info
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
	sqlstr := `insert into user_sim(uid,sid,sim_provider,sim_no,iccid,sim_byte,sim_available_byte,bind_ts,service_end_ts,service_duration,remark,created,updated)
	values(%d,%d,'%s','%s','%s',%d,%d,%d,%d,%d,'%s',%d,%d)
	returning uid`
	sqlstr = fmt.Sprintf(sqlstr, i.Uid, i.Sid, i.SimProvider, i.SimNo, i.Iccid, i.SimByte, i.SimAvailableByte, i.BindTs, i.ServiceEndTs, i.ServiceDuration, i.Remark, i.Created, i.Updated)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", i, sqlstr, err)
		return
	}
	return
}

func (o object) Update(i *Info) (r int) {
	now := time.Now()
	sqlstr := fmt.Sprintf("update user_sim set updated = %d,", now.Unix())
	if i.Uid <= 0 && i.Sid < 0 {
		o.Log.Error("invalid param")
		return
	}
	if i.ServiceEndTs > 0 {
		sqlstr += fmt.Sprintf("service_end_ts = %d,", i.ServiceEndTs)
	}
	if i.ServiceDuration > 0 {
		sqlstr += fmt.Sprintf("service_duration = %d,", i.ServiceDuration)
	}
	if i.SimByte > 0 {
		sqlstr += fmt.Sprintf("sim_byte = %d,", i.SimByte)
	}
	if i.SimAvailableByte > 0 {
		sqlstr += fmt.Sprintf("sim_available_byte = %d,", i.SimAvailableByte)
	}
	if len(i.Remark) > 0 {
		sqlstr += fmt.Sprintf("remark = '%s',", i.Remark)
	}
	sqlstr = strings.TrimRight(sqlstr, ",")
	sqlstr += fmt.Sprintf(" where uid = %d and sid = %d returning uid", i.Uid, i.Sid)
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
	coalesce(sim_provider,'') as sim_provider,
	coalesce(sim_no,'') as sim_no,
	coalesce(iccid,'') as iccid,
	coalesce(sim_byte,0) as sim_byte,
	coalesce(sim_available_byte,0) as sim_available_byte,
	coalesce(bind_ts,0) as bind_ts,
	coalesce(service_end_ts,0) as service_end_ts,
	coalesce(service_duration,0) as service_duration,
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
		err = rows.Scan(&r.Uid, &r.Sid, &r.SimProvider, &r.SimNo, &r.Iccid, &r.SimByte, &r.SimAvailableByte, &r.BindTs, &r.ServiceEndTs, &r.ServiceDuration, &r.Remark, &r.Created, &r.Updated)
		if err != nil {
			o.Log.Errorf("param=%+v,sql=%s,err=%v", uid, sqlstr, err)
			continue
		}
		list = append(list, r)
	}
	return
}

func (o object) queryInfo(uid, sid int) *Info {
	r := new(Info)
	var qstr string
	if uid > 0 {
		qstr += fmt.Sprintf(" and uid = %d ", uid)
	}
	if sid > 0 {
		qstr += fmt.Sprintf(" and sid = %d ", sid)
	}
	sqlstr := `
	select
	coalesce(uid,0) as uid,
	coalesce(sid,0) as sid,
	coalesce(sim_provider,'') as sim_provider,
	coalesce(sim_no,'') as sim_no,
	coalesce(iccid,'') as iccid,
	coalesce(sim_byte,0) as sim_byte,
	coalesce(sim_available_byte,0) as sim_available_byte,
	coalesce(bind_ts,0) as bind_ts,
	coalesce(service_end_ts,0) as service_end_ts,
	coalesce(service_duration,0) as service_duration,
	coalesce(remark,'') as remark,
	coalesce(created,0) as created,
	coalesce(updated,0) as updated
	from user_sim
	where 1=1
	`
	sqlstr += qstr
	// app.Log.Debug(sqlstr)
	err := o.DbRo.QueryRow(sqlstr).Scan(&r.Uid, &r.Sid, &r.SimProvider, &r.SimNo, &r.Iccid, &r.SimByte, &r.SimAvailableByte, &r.BindTs, &r.ServiceEndTs, &r.ServiceDuration, &r.Remark, &r.Created, &r.Updated)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", uid, sqlstr, err)
		return nil
	}
	return r
}

func (o object) QueryInfo(uid, sid int) *Info {
	return o.queryInfo(uid, sid)
}

func (o object) QueryInfoBySid(sid int) *Info {
	return o.queryInfo(0, sid)
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
