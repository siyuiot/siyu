package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/minappCode2PhoneNum"
	"github.com/siyuiot/siyu/app/user/internal/minappCode2Session"
	"github.com/siyuiot/siyu/app/user/internal/minappDataDecrypt"
	"github.com/siyuiot/siyu/app/user/internal/user"
	"github.com/siyuiot/siyu/app/user/internal/userToken"
	"github.com/siyuiot/siyu/app/user/internal/wechatAccessToken"
	"github.com/siyuiot/siyu/modules/qstate"
	"github.com/siyuiot/siyu/pkg/qmd5"
	"github.com/siyuiot/siyu/pkg/quuid"
)

func (t *HttpServer) LoginMinapp(c *gin.Context) {
	type (
		Data struct {
			Token   string         `json:"token"`
			Expires int64          `json:"expires"`
			OpenId  string         `json:"openId"`
			User    *user.UserInfo `json:"user"`
		}
		Req struct {
			Code   string                   `json:"code" binding:"required"`
			Minapp minappDataDecrypt.Minapp `json:"minapp"`
		}
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			Data      Data   `json:"data"`
		}
	)
	var req Req
	var ret Ret
	defer func() {
		c.JSON(http.StatusOK, ret)
	}()
	err := c.ShouldBindJSON(&req)
	if err != nil {
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
		return
	}
	wat := wechatAccessToken.Instance().GetFromDbOrRemote(req.Minapp.AppId)
	if wat == nil {
		ret.State, ret.StateInfo = qstate.StateFailed, "accessToken is nil"
		return
	}
	ms := minappCode2Session.GET(wat.AppId, wat.Secret, req.Minapp.Code)
	if ms == nil {
		ret.State, ret.StateInfo = qstate.StateFailed, "session nil,may be invalid code"
		return
	}
	data := Data{OpenId: ms.OpenID}
	// phoneInfo, err := minappDataDecrypt.VerifXcxUserPnInfo(ms.SessionKey, req.Minapp.Iv, req.Minapp.EncryptedData, req.Minapp.AppId)
	phoneInfo := minappCode2PhoneNum.POST(wat.AccessToken, req.Code)
	if phoneInfo == nil {
		ret.State, ret.StateInfo = qstate.StateFailed, "phoneInfo is nil,may be invalid code"
		return
	}
	if len(phoneInfo.PhoneNumber) != 11 {
		ret.State, ret.StateInfo = qstate.StateFailed, fmt.Sprintf("phoneNumber:%s,may be invalid code", phoneInfo.PhoneNumber)
		return
	}
	app := "com.siyu.iot"
	now := time.Now().UTC()
	u := user.Instance().QueryInfoByPn(phoneInfo.PhoneNumber, phoneInfo.CountryCode, app)
	if u == nil {
		u = &user.UserInfo{
			PhoneNum:    phoneInfo.PhoneNumber,
			PhoneArea:   phoneInfo.CountryCode,
			App:         app,
			PwdSalt:     quuid.New(),
			NickName:    fmt.Sprintf("%s%d", strings.Split(app, ".")[1], now.Unix()),
			Completion:  30,
			CreatedTime: now,
			UpdatedTime: now,
			RegType:     1,
		}
		user.Instance().Insert(u)
	}
	data.Token = qmd5.QMD5String([]byte(fmt.Sprintf("foo%dbar", now.Unix())))
	data.User = u
	data.Expires = now.Add(7200 * time.Second).Unix()
	ruid := userToken.Instance().Upsert(&userToken.Info{Uid: u.UserId, Ts: now.Unix(), Token: data.Token, Expires: data.Expires})
	if ruid == false {
		ret.State, ret.StateInfo = State(qstate.StateFailed)
		return
	}
	ret.Data = data
	ret.State, ret.StateInfo = State(qstate.StateOk)
}
