package sim

import (
	"database/sql"
	"fmt"

	"github.com/siyuiot/siyu/pkg/qlog"
)

var this *object

type Object interface {
	QueryInfoByIccid(id string) *Info
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

func (o object) queryInfo(id int, simNo, imsi, iccid string) *Info {
	r := new(Info)
	var qstr string
	switch {
	case len(simNo) > 0:
		qstr += fmt.Sprintf(" AND sim_no = '%s'", simNo)
	case len(imsi) > 0:
		qstr += fmt.Sprintf(" AND imsi = '%s'", imsi)
	case len(iccid) > 0:
		qstr += fmt.Sprintf(" AND iccid = '%s'", iccid)
	default:
		o.Log.Error("invalid param")
		return nil
	}
	sqlstr := `
	select
	coalesce(id,0) as id,
	coalesce(sim_no,'') as sim_no,
	coalesce(imsi,'') as imsi,
	coalesce(iccid,'') as iccid,
	coalesce(Remark,'') as Remark
	from sim
	where 1=1
	`
	sqlstr += qstr
	// app.Log.Debug(sqlstr)
	err := o.DbRo.QueryRow(sqlstr).Scan(&r.Id, &r.SimNo, &r.Imsi, &r.Iccid, &r.Remark)
	if err != nil {
		o.Log.Errorf("param=%s,sql=%s,err=%v", id, sqlstr, err)
		return nil
	}
	return r
}

func (o object) QueryInfoByIccid(iccid string) *Info {
	return o.queryInfo(0, "", "", iccid)
}
