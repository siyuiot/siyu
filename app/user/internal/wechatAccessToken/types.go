package wechatAccessToken

type Info struct {
	AppId       string
	Secret      string
	AccessToken string
	ExpiresIn   int
	ExpiresAt   int64
	Remark      string
	Created     int64
	Updated     int64
}
