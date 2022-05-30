package minappCode2Session

import (
	"fmt"

	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/pkg/qhttp"
)

// code换取登录凭证校验。 每个 code 只能使用一次，code的有效期为5min

// 官方文档https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html

// curl -H "Accept: application/json" -H "Content-type: application/json" -X GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code

func GET(appId, secret, code string) *Info {
	type Out struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
		Info
	}
	if len(secret) <= 0 {
		app.Log.Errorf("secret:%s", secret)
		return nil
	}
	if len(code) <= 0 {
		app.Log.Errorf("code:%s", code)
		return nil
	}
	var out Out
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, secret, code)
	resp, err := qhttp.GetJSON(url, &out)
	if err != nil {
		app.Log.Errorf("err:%s", err)
		return nil
	}
	if len(out.Info.SessionKey) <= 0 {
		app.Log.Errorf("SessionKey:%s,out:%#+v,err:%s", out.Info.SessionKey, out, resp)
		return nil
	}
	return &out.Info
}
