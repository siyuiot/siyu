package minapp

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"

	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/pkg/qhttp"
	"github.com/spf13/viper"
)

var (
	ErrAppIDNotMatch       = errors.New("app id not match")
	ErrInvalidBlockSize    = errors.New("invalid block size")
	ErrInvalidPKCS7Data    = errors.New("invalid PKCS7 data")
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type WXBizDataCrypt struct {
	appID, sessionKey string
}

func NewWXBizDataCrypt(appID, sessionKey string) *WXBizDataCrypt {
	return &WXBizDataCrypt{
		appID:      appID,
		sessionKey: sessionKey,
	}
}

// pkcs7Unpad returns slice of the original data without padding
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(data)%blockSize != 0 || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	c := data[len(data)-1]
	n := int(c)
	if n == 0 || n > len(data) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if data[len(data)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return data[:len(data)-n], nil
}

func (w *WXBizDataCrypt) Decrypt(encryptedData, iv string) ([]byte, error) {
	aesKey, err := base64.StdEncoding.DecodeString(w.sessionKey)
	if err != nil {
		return nil, err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(cipherText, cipherText)
	return pkcs7Unpad(cipherText, block.BlockSize())
}

type SessionOpenIDStr struct {
	SessionKey string `json:"session_key"`
	Expires    int    `json:"expires_in"`
	XcxOpenID  string `json:"openid"`
}

func VerifyXcxCode(code, appid string) (openid, sessionkey string, expire int, err error) {
	var clusterinfo = url.Values{}
	clusterinfo.Add("appid", appid)
	secret := viper.GetString("minapp." + appid)
	if secret == "" {
		app.Log.Errorf("appid:%v 未配置 secret", appid)
		err = errors.New("未配置 secret")
		return
	}
	clusterinfo.Add("secret", secret)
	clusterinfo.Add("js_code", code)
	clusterinfo.Add("grant_type", "authorization_code")
	sos := &SessionOpenIDStr{}
	reqUrl := "https://api.weixin.qq.com/sns/jscode2session?" + clusterinfo.Encode()
	body, err := qhttp.GetJSON(reqUrl, sos)
	app.Log.Debugf(string(body))
	if err != nil {
		app.Log.Error("json Unmarshal content = ", string(body), " err = ", err)
		return "", "", 0, err
	}
	return sos.XcxOpenID, sos.SessionKey, sos.Expires, nil
}

func VerifXcxUserInfo(sessionkey, iv, data, appid string) (unionid string, err error) {
	pc := NewWXBizDataCrypt(appid, sessionkey)
	b, err := pc.Decrypt(data, iv)
	if err != nil {
		app.Log.Error("VerifXcxEncryptedData err = ", err)
		return "", err
	}
	type WechatUserInfo struct {
		OpenID    string `json:"openId"`
		UnionID   string `json:"unionId"`
		NickName  string `json:"nickName"`
		Gender    int    `json:"gender"`
		City      string `json:"city"`
		Province  string `json:"province"`
		Country   string `json:"country"`
		AvatarURL string `json:"avatarUrl"`
		Language  string `json:"language"`
		Watermark struct {
			Timestamp int64  `json:"timestamp"`
			AppID     string `json:"appid"`
		} `json:"watermark"`
	}
	var ui WechatUserInfo
	err = json.Unmarshal(b, &ui)
	if err != nil {
		return "", err
	}
	if ui.Watermark.AppID != appid {
		return "", errors.New("app id not match")
	}
	return ui.UnionID, nil
}

func VerifXcxUserPnInfo(sessionkey, iv, data, appid string) (pn string, err error) {
	pc := NewWXBizDataCrypt(appid, sessionkey)
	b, err := pc.Decrypt(data, iv)
	if err != nil {
		app.Log.Error("VerifXcxEncryptedData err = ", err)
		return "", err
	}
	type WechatUserPnInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
		Watermark       struct {
			Timestamp int64  `json:"timestamp"`
			AppID     string `json:"appid"`
		} `json:"watermark"`
	}
	var ui WechatUserPnInfo
	err = json.Unmarshal(b, &ui)
	if err != nil {
		return "", err
	}
	if ui.Watermark.AppID != appid {
		return "", errors.New("app id not match")
	}
	return ui.PurePhoneNumber, nil
}
