package user

// const (
// 	sqlQueryusercheckinstat = `
// select to_char(last_checkin_date,'YYYY-MM-DD'), last_checkin_time, consec_days
//   from user_checkin_stat
//  where user_id = $1`
// 	sqlUpdateusercheckinstat = `
// with a as (
// update user_checkin_stat
//    set last_checkin_date = $2::date
//       ,last_checkin_time = $3
//  where user_id = $1
// returning user_id
// )
// insert into user_checkin_stat(user_id, last_checkin_date, last_checkin_time, consec_days)
// select $1, $2::date, $3, 1
//  where not exists (select 1 from a)`

// 	sqlUpsertLoginToken = `INSERT INTO user_login_token(uid, ts, token, expire, des)
// VALUES ($1, $2, $3, $4, $5)
// ON CONFLICT (uid)
// DO
//    UPDATE SET ts=$2, token=$3, expire=$4, des=$5;`
// 	sqlQueryUserThirdMap = `select user_id, partner, partner_openid, partner_uid, partner_nick_name, binded_time, app from user_third_map where user_id=$1`

// 	sqlQueryUserBikeMap = `select user_id, bike_id, is_master, bike_name, binded_time, status, security_mode, security_custom, created_time, display_order, bike_seq from user_bike_map where user_id=$1`

// 	sqlQueryUserWxMap = `select user_id, unionid, app_openid, pub_openid, xcx_openid, created_time from user_wx_map where user_id=$1`
// )
