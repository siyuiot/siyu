package wechat

import "github.com/siyuiot/siyu/modules/qstate"

type (
	InCode2Session struct {
		Appid string `json:"appid" binding:"required"`
		Code  string `json:"code"`
		Token string `json:"token"`
	}
	OutCode2Session struct {
		qstate.CommonResponse
		Data Code2SessionData `json:"data"`
	}
	Code2SessionData struct {
		OpenID     string `json:"openid"`  // 用户唯一标识
		UnionID    string `json:"unionid"` // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
		SessionKey string `json:"sessionkey"`
		Token      string `json:"token"`
	}
)
