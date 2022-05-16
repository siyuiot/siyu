package user

import "time"

// type UserInfo struct {
// 	User_Id     int32    `json:"user_id"`
// 	Phone_Num   string   `json:"phone_num"`
// 	Account     string   `json:"account"`
// 	Email       string   `json:"email"`
// 	Nick_Name   string   `json:"nick_name"`
// 	Icon        string   `json:"icon"`
// 	Real_Name   string   `json:"real_name"`
// 	Gender      int16    `json:"gender"`
// 	Birthday    string   `json:"birthday"`
// 	ID_No       string   `json:"id_no"`
// 	Location    string   `json:"location"`
// 	Reg_Time    string   `json:"red_time"`
// 	Pwd_Status  int16    `json:"pwd_status"`
// 	Per_Sign    string   `json:"per_sign"`
// 	Completion  int16    `json:"completion"`
// 	Home        Position `json:"home"`
// 	Mile_Remind int      `json:"mile_remind"`
// 	MileType    int      `json:"mile_type"`
// }

// type Info struct {
// 	UserId       int       `json:"user_id"`
// 	PhoneNum     string    `json:"phone_num"`
// 	PhoneArea    string    `json:"phone_area"`
// 	Account      string    `json:"account"`
// 	Email        string    `json:"email"`
// 	Password     string    `json:"-"`
// 	PwdSalt      string    `json:"-"`
// 	NickName     string    `json:"nick_name"`
// 	RealName     string    `json:"real_name"`
// 	Gender       int16     `json:"gender"`
// 	Birthday     time.Time `json:"birthday"` //TODO
// 	IdNo         string    `json:"id_no"`
// 	Icon         string    `json:"icon"`
// 	Location     string    `json:"location"`
// 	Reg_Time     string    `json:"red_time"`
// 	PerSign      string    `json:"per_sign"`
// 	Completion   int       `json:"completion"`
// 	Home         string    `json:"home"`
// 	RegType      int16
// 	Tz           string
// 	App          string
// 	GeneralSetup string
// 	CreatedTime  time.Time
// 	UpdatedTime  time.Time
// }

type UserInfo struct {
	UserId       int       `json:"user_id"`
	PhoneNum     string    `json:"phone_num"`
	PhoneArea    string    `json:"phone_area"`
	Account      string    `json:"account"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	PwdSalt      string    `json:"-"`
	NickName     string    `json:"nick_name"`
	RealName     string    `json:"real_name"`
	Gender       int16     `json:"gender"`
	Birthday     string    `json:"birthday"`
	IdNo         string    `json:"id_no"`
	Icon         string    `json:"icon"`
	Location     string    `json:"location"`
	Reg_Time     string    `json:"red_time"`
	PerSign      string    `json:"per_sign"`
	Completion   int       `json:"completion"`
	Home         string    `json:"-"`
	RegType      int16     `json:"-"`
	Tz           string    `json:"-"`
	App          string    `json:"-"`
	GeneralSetup string    `json:"-"`
	CreatedTime  time.Time `json:"-"`
	UpdatedTime  time.Time `json:"-"`
}

type LoginToken struct {
	Uid    int
	Token  string
	Ts     int64
	Des    string
	Expire int64
}
