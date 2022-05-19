package minappCode2PhoneNum

import (
	"fmt"

	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/pkg/qhttp"
)

// code换取用户手机号。 每个 code 只能使用一次，code的有效期为5min

// 官方文档https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/phonenumber/phonenumber.getPhoneNumber.html

// curl -H "Accept: application/json" -H "Content-type: application/json" -X POST -d '{"code": "e31968a7f94cc5ee25fafc2aef2773f0bb8c3937b22520eb8ee345274d00c144"}' https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=ACCESS_TOKEN&
// {
//     "errcode":0,
//     "errmsg":"ok",
//     "phone_info": {
//         "phoneNumber":"xxxxxx",
//         "purePhoneNumber": "xxxxxx",
//         "countryCode": 86,
//         "watermark": {
//             "timestamp": 1637744274,
//             "appid": "xxxx"
//         }
//     }
// }

func POST(accessToken, code string) *Info {
	type In struct {
		Code string `json:"code"`
	}
	type Out struct {
		Errcode   int    `json:"errcode"`
		Errmsg    string `json:"errmsg"`
		PhoneInfo Info   `json:"phone_info"`
	}
	if len(accessToken) <= 0 {
		app.Log.Errorf("accessToken:%s", accessToken)
		return nil
	}
	if len(code) <= 0 {
		app.Log.Errorf("code:%s", code)
		return nil
	}
	var out Out
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s&", accessToken)
	// header := map[string]string{
	// 	// "Accept":       "application/json",
	// 	// "Content-type": "application/json",
	// }
	var in = In{Code: code}
	// app.Log.Debugf("url:%s,header:%#+v,in:%#+v", url, header, in)
	// b, err := json.Marshal(in)
	// app.Log.Debug(string(b))
	resp, err := qhttp.PostJSON(url, nil, in, &out)
	if err != nil {
		app.Log.Errorf("err:%s", err)
		return nil
	}
	if len(out.PhoneInfo.PhoneNumber) <= 0 {
		app.Log.Errorf("phoneNumber:%s,out:%#+v,err:%v", out.PhoneInfo.PhoneNumber, out, resp)
		return nil
	}
	return &out.PhoneInfo
}
