package httpserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/minapp"
	"github.com/siyuiot/siyu/app/user/internal/user"
	"github.com/siyuiot/siyu/modules/qstate"
	"github.com/siyuiot/siyu/pkg/qmd5"
	"github.com/siyuiot/siyu/pkg/quuid"
)

func (t *HttpServer) LoginMinapp(c *gin.Context) {
	type (
		Data struct {
			Token   string         `json:"token"`
			Expires int64          `json:"expires"`
			User    *user.UserInfo `json:"user"`
		}
		Req struct {
			Minapp minapp.Minapp `json:"minapp"`
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
		ret.State, ret.StateInfo = State(qstate.StateInvalidParameter)
		return
	}
	if len(req.Minapp.AppId) == 0 || len(req.Minapp.Iv) == 0 || len(req.Minapp.EncryptedData) == 0 {
		app.Log.Error("invalid parametes")
		ret.State, ret.StateInfo = State(qstate.StateInvalidParameter)
		return
	}
	// out, err := wechat.New(wechat.ClientOption{
	// 	RequestId: quuid.New(),
	// }).Code2Session(wechat.InCode2Session{
	// 	Appid: req.Minapp.AppId,
	// 	Code:  req.Minapp.Code,
	// 	Token: req.Minapp.TokenMp,
	// })
	// if err != nil {
	// 	app.Log.Error(err)
	// 	ret.State, ret.StateInfo = State(qstate.StateFailed)
	// 	return
	// }
	// pn, err := minapp.VerifXcxUserPnInfo(out.SessionKey, req.Minapp.Iv, req.Minapp.EncryptedData, req.Minapp.AppId)
	// if err != nil {
	// 	app.Log.Error("invalid parametes", err)
	// 	ret.State, ret.StateInfo = State(qstate.StateInvalidParameter)
	// 	return
	// }
	pn := "18616854987"
	phoneAreaMainLand := "86"
	app := "com.siyu.iot"
	now := time.Now()
	u := t.UserEntry.QueryInfoByPn(pn, phoneAreaMainLand, app)
	if u == nil {
		u = &user.UserInfo{
			PhoneNum:    pn,
			PhoneArea:   phoneAreaMainLand,
			App:         app,
			PwdSalt:     quuid.New(),
			NickName:    fmt.Sprintf("%s%s", strings.Split(app, ".")[1], now.Format("0102150405")),
			Completion:  30,
			CreatedTime: now,
			UpdatedTime: now,
			RegType:     1,
		}
		t.UserEntry.Insert(u)
	}
	data := Data{
		Token: qmd5.QMD5String([]byte(fmt.Sprintf("foo%dbar", now.Unix()))),
		User:  u,
	}
	data.Expires = now.Add(7200 * time.Second).Unix()
	ruid := t.UserEntry.UpsertLoginToken(user.LoginToken{Uid: u.UserId, Ts: now.Unix(), Token: data.Token, Expires: data.Expires})
	if ruid == false {
		ret.State, ret.StateInfo = State(qstate.StateFailed)
		return
	}
	ret.Data = data
	ret.State, ret.StateInfo = State(qstate.StateOk)
}
