package user

//用户信息

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Object interface {
	Insert(info *UserInfo) int
	Update(info *UserInfo) int
	UpdatePnIsNull(uid int) int
	UpdateEmailIsNull(uid int) int
	UpdatePasswd(uid int, pw string) int
	QueryInfo(uid int, pn, pnArea, app, email string) *UserInfo
	QueryInfoByUid(uid int) *UserInfo
	QueryInfoByPn(pn, pnArea, app string) *UserInfo
	QueryInfoByEmail(email, app string) *UserInfo
	Delete(uid int) (r int)
	UpsertLoginToken(info LoginToken) (ok bool)
}

type Option struct {
	Log *logrus.Entry
	Db  *sql.DB
}

type object struct {
	Option
}

func New(o Option) Object {
	return &object{Option: o}
}

func (o object) Insert(u *UserInfo) (r int) {
	sqlstr := `insert into users(created_time,updated_time,phone_num,email,phone_area,app,nick_name,pwd_salt,reg_type,password)
	values('%s','%s',%s,'%s','%s','%s','%s','%s','%d','%s')
	returning user_id`
	if len(u.PhoneNum) <= 0 { //数据库phone_num有唯一约束，插入NULL
		u.PhoneNum = "NULL"
	} else {
		u.PhoneNum = fmt.Sprintf("'%s'", u.PhoneNum)
	}
	sqlstr = fmt.Sprintf(sqlstr, u.CreatedTime.Format(time.RFC3339), u.UpdatedTime.Format(time.RFC3339), u.PhoneNum, u.Email, u.PhoneArea, u.App, u.NickName, u.PwdSalt, u.RegType, u.Password)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%+v,sql=%s,err=%v", u, sqlstr, err)
		return
	}
	u.UserId = r
	return
}

func (o object) Update(u *UserInfo) (r int) {
	sqlstr := fmt.Sprintf("update users set updated_time = '%s',", time.Now().Format(time.RFC3339))
	if u.UserId <= 0 {
		o.Log.Error("invalid param")
		return
	}
	if len(u.PhoneNum) > 0 {
		sqlstr += fmt.Sprintf("phone_num = '%s',", u.PhoneNum)
	}
	if len(u.PhoneArea) > 0 {
		sqlstr += fmt.Sprintf("phone_area = '%s',", u.PhoneArea)
	}
	if len(u.Email) > 0 {
		sqlstr += fmt.Sprintf("email = '%s',", u.Email)
	}
	if len(u.NickName) > 0 {
		sqlstr += fmt.Sprintf("nick_name = '%s',", u.NickName)
	}
	if len(u.RealName) > 0 {
		sqlstr += fmt.Sprintf("real_name = '%s',", u.RealName)
	}
	if u.Gender > 0 {
		sqlstr += fmt.Sprintf("gender = %d,", u.Gender)
	}
	if b, _ := time.Parse("2006-01-02", u.Birthday); (b != time.Time{}) {
		sqlstr += fmt.Sprintf("birthday = '%s',", u.Birthday)
	}
	if len(u.IdNo) > 0 {
		sqlstr += fmt.Sprintf("id_no = '%s',", u.IdNo)
	}
	if len(u.Icon) > 0 {
		sqlstr += fmt.Sprintf("icon = '%s',", u.Icon)
	}
	if len(u.PerSign) > 0 {
		sqlstr += fmt.Sprintf("per_sign = '%s',", u.PerSign)
	}
	if u.Completion > 0 {
		sqlstr += fmt.Sprintf("completion = %d,", u.Completion)
	}
	if len(u.Home) > 0 {
		sqlstr += fmt.Sprintf("home = '%s',", u.Home)
	}
	if len(u.Location) > 0 {
		sqlstr += fmt.Sprintf("location = '%s',", u.Location)
	}
	if len(u.GeneralSetup) > 0 {
		sqlstr += fmt.Sprintf("general_setup = '%s',", u.GeneralSetup)
	}
	sqlstr = strings.TrimRight(sqlstr, ",")
	sqlstr += fmt.Sprintf(" where user_id = %d returning user_id", u.UserId)
	err := o.Db.QueryRow(sqlstr).Scan(&r)
	if err != nil {
		o.Log.Errorf("user=%+v,sql=%s,err=%v", u, sqlstr, err)
		return
	}
	return
}

func (o object) UpdatePnIsNull(uid int) (r int) {
	sqlstr := `
	update users set 
	phone_num = null
	where user_id = $1
	returning user_id
	`
	err := o.Db.QueryRow(sqlstr, uid).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", uid, sqlstr, err)
		return
	}
	return
}

func (o object) UpdateEmailIsNull(uid int) (r int) {
	sqlstr := `
	update users set 
	email = null
	where user_id = $1
	returning user_id
	`
	err := o.Db.QueryRow(sqlstr, uid).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", uid, sqlstr, err)
		return
	}
	return
}

func (o object) UpdatePasswd(uid int, pw string) (r int) {
	sqlstr := `
	update users
	set password = $2
	where user_id = $1
	returning user_id
	`
	err := o.Db.QueryRow(sqlstr, uid, pw).Scan(&r)
	if err != nil {
		o.Log.Errorf("param=%d %s,sql=%s,err=%v", uid, pw, sqlstr, err)
		return
	}
	return
}

func (o object) QueryInfoByUid(uid int) *UserInfo {
	return o.queryInfo(uid, "", "", "", "")
}

func (o object) QueryInfoByPn(pn, pnArea, app string) *UserInfo {
	return o.queryInfo(0, pn, pnArea, app, "")
}

func (o object) QueryInfoByEmail(email, app string) *UserInfo {
	return o.queryInfo(0, "", "", app, email)
}

func (o object) QueryInfo(uid int, pn, pnArea, app, email string) *UserInfo {
	return o.queryInfo(uid, pn, pnArea, app, email)
}

func (o object) queryInfo(uid int, pn, pnArea, app, email string) *UserInfo {
	u := new(UserInfo)
	var qstr string
	switch {
	case uid > 0:
		qstr = fmt.Sprintf(" and u.user_id = %d", uid)
	case len(pn) > 0:
		qstr = fmt.Sprintf(" and u.phone_num = '%s' and u.phone_area = '%s' and u.app= '%s' ", pn, pnArea, app)
	case len(email) > 0:
		qstr = fmt.Sprintf(" and u.email = '%s' and u.app= '%s' ", email, app)
	default:
		o.Log.Debug("invalid parametes")
		return nil
	}
	sqlstr := `
	select 
	u.user_id as user_id,
	coalesce(u.phone_num,''),
	coalesce(u.email,''),
	u.nick_name,
	coalesce(u.real_name,''),
	coalesce(u.gender,0),
	coalesce(u.birthday,'0001-01-01'::date),
	coalesce(u.id_no,''),
	coalesce(u.icon,''),
	coalesce(u.location,''),
	u.created_time,
	coalesce(u.password,''),
	coalesce(u.per_sign,''),
	coalesce(u.app,''),
	u.completion,
	u.home,
	u.general_setup,
	u.tz,
	u.pwd_salt,
	coalesce(u.phone_area,'')
  	from users u
 	where 1=1
	`
	sqlstr += qstr
	var birthDay time.Time
	err := o.Db.QueryRow(sqlstr).Scan(&u.UserId, &u.PhoneNum, &u.Email, &u.NickName, &u.RealName, &u.Gender, &birthDay, &u.IdNo, &u.Icon, &u.Location, &u.CreatedTime, &u.Password, &u.PerSign,
		&u.App, &u.Completion, &u.Home, &u.GeneralSetup, &u.Tz, &u.PwdSalt, &u.PhoneArea)
	if err != nil {
		o.Log.Errorf("param=%d,%s,sql=%s,err=%v", uid, pn, sqlstr, err)
		return nil
	}
	u.Birthday = birthDay.Format("2006-01-02")
	//tmp := user.GeneralSetupStr{}
	//err = json.Unmarshal([]byte(u.GeneralSetup), &tmp)
	//if err != nil {
	//	o.Log.Warnf("Unmarshal json error=%v", err)
	//}
	//u.GenSetStr = tmp
	return u
}

func (o object) queryList(where string, page, size int) (list []UserInfo, total int) {
	list = make([]UserInfo, 0)
	if page <= 0 {
		page = 1
	}
	if page <= 0 {
		size = 20
	}
	sqlCount := `select count(user_id) from users where 1=1` + where
	err := o.Db.QueryRow(sqlCount).Scan(&total)
	if err != nil {
		o.Log.Errorf("sqlCount=%s,error=%v", sqlCount, err)
	}
	sqlstr := `
	select
	user_id,
	coalesce(phone_num,''),
	coalesce(account,''),
	coalesce(email,''),
	coalesce(nick_name,''),
	coalesce(real_name,''),
	coalesce(gender,0),
	coalesce(birthday,'0001-01-01'::date),
	coalesce(id_no,''),
	coalesce(icon,''),
	coalesce(password,''),
	coalesce(location,''),
	created_time,
	updated_time,
	coalesce(pwd_salt,''),
	coalesce(per_sign,''),
	coalesce(completion,0),
	coalesce(home,'{}'::jsonb),
	coalesce(reg_type,0),
	coalesce(tz,''),
	coalesce(phone_area,''),
	coalesce(app,''),
	coalesce(general_setup,'{}'::jsonb)
	from users  where 1=1`
	where = sqlstr + where + fmt.Sprintf(" limit %d offset %d", size, (page-1)*size)
	rows, err := o.Db.Query(where)
	if err != nil {
		o.Log.Errorf("sqlstr=%s,error=%v", sqlstr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tmp := UserInfo{}
		var birthDay time.Time
		err := rows.Scan(&tmp.UserId, &tmp.PhoneNum, &tmp.Account, &tmp.Email, &tmp.NickName, &tmp.RealName, &tmp.Gender, &tmp.Birthday,
			&tmp.IdNo, &tmp.Icon, &tmp.Password, &tmp.Location, &tmp.CreatedTime, &tmp.UpdatedTime, &tmp.PwdSalt, &tmp.PerSign,
			&tmp.Completion, &tmp.Home, &tmp.RegType, &tmp.Tz, &tmp.PhoneArea, &tmp.App, &tmp.GeneralSetup)
		if err != nil {
			o.Log.Error(err)
			continue
		}
		tmp.Birthday = birthDay.Format("2006-01-02")
		list = append(list, tmp)
	}

	return
}

func (o object) Delete(uid int) (r int) {
	sqlstr := `
	do $$
	declare uid integer;
	begin
	uid := %d;
	delete from user_wx_map where user_id = uid;
	delete from user_third_map where user_id = uid;
	delete from user_bike_map where user_id = uid;
	update bike set owner_id = null where owner_id = uid;
	delete from users where user_id = uid;
	end;
	$$
	`
	sqlstr = fmt.Sprintf(sqlstr, uid)
	_, err := o.Db.Exec(sqlstr)
	if err != nil {
		o.Log.Errorf("param=%d,sql=%s,err=%v", uid, sqlstr, err)
		return
	}
	return uid
}

func (o *Option) UpsertLoginToken(info LoginToken) (ok bool) {
	_, err := o.Db.Exec(`INSERT INTO user_login_token(uid, ts, token, expire, des)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (uid)
	DO UPDATE SET ts=$2, token=$3, expire=$4, des=$5;`, info.Uid, info.Ts, info.Token, info.Expires, info.Des)
	if err != nil {
		o.Log.Errorf("upsert login_token err(%v) ", err)
		return false
	}
	return true
}
