package httpserver

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/app/user/internal/app"
	"github.com/siyuiot/siyu/app/user/internal/sim"
	"github.com/siyuiot/siyu/app/user/internal/userSim"
	"github.com/siyuiot/siyu/app/user/internal/userSimLog"
	"github.com/siyuiot/siyu/modules/qstate"
)

func (t *HttpServer) UserSimBind(c *gin.Context) {
	type (
		Req struct {
			IccId string `json:"iccid"`
			Pku   string `json:"pku"`
		}
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			CustInfo  string `json:"custInfo,omitempty"` //自定义描述
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
	ruid, _ := strconv.Atoi(c.Request.Header.Get("ruid"))
	// 查询sim信息
	info := sim.Instance().QueryInfoByIccid(req.IccId)
	if info == nil {
		err := fmt.Errorf("iccid not found")
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
		return
	}
	app.Log.Debugf("%#+v", info)
	// 查询用户sim绑定关系
	us := userSim.Instance().QueryInfo(ruid, info.Id)
	if us != nil {
		err := fmt.Errorf("sid:%d already bind uid:%d", info.Id, us.Uid)
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
		return
	}
	// 绑定用户和sim卡
	now := time.Now()
	r := userSim.Instance().Insert(&userSim.Info{
		Uid:              ruid,
		Sid:              info.Id,
		SimNo:            info.SimNo,
		SimByte:          307200,
		SimAvailableByte: 307200,
		BindTs:           now.Unix(),
		ServiceEndTs:     now.Unix() + 86400, // 绑卡送1天体验
		Remark:           fmt.Sprintf("用户:%d,卡号:%s,iccid:%s", ruid, info.SimNo, info.Iccid),
	})
	if r <= 0 {
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, "bind error"
		return
	}
	// 添加绑定记录
	userSimLog.Instance().Insert(&userSimLog.Info{
		Ts:       now.Unix(),
		Uid:      ruid,
		Sid:      info.Id,
		PhoneNum: strconv.Itoa(ruid),
		SimNo:    info.SimNo,
		Imsi:     info.Imsi,
		IccId:    info.Iccid,
		Remark:   "UserSimBind",
	})
	ret.State, ret.StateInfo = State(qstate.StateOk)
}

func (t *HttpServer) UserSimList(c *gin.Context) {
	type (
		Data struct {
			Sid   int    `json:"sid"`
			SimNo string `json:"simNo"`
		}
		Req struct {
		}
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			CustInfo  string `json:"custInfo,omitempty"` //自定义描述
			Data      []Data `json:"data"`
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
	ruid, _ := strconv.Atoi(c.Request.Header.Get("ruid"))
	// 查询用户sim绑定关系
	_, list := userSim.Instance().QueryList(ruid)
	for _, v := range list {
		ret.Data = append(ret.Data, Data{Sid: v.Sid, SimNo: v.SimNo})
	}
	ret.State, ret.StateInfo = State(qstate.StateOk)
}

func (t *HttpServer) UserSimInfo(c *gin.Context) {
	type (
		Data struct {
			Sid              int    `json:"sid"`
			SimNo            string `json:"simNo"`
			BindTs           int64  `json:"bindTs"`           // 绑定时间戳
			ServiceEndTs     int64  `json:"serviceEndTs"`     // 服务结束时间戳
			SimByte          int    `json:"simByte"`          // 卡流量
			SimAvailableByte int    `json:"simAvailableByte"` // 可用流量
		}
		Req struct {
			Sid int `json:"sid"`
		}
		Ret struct {
			State     int    `json:"state"`
			StateInfo string `json:"stateInfo"`
			CustInfo  string `json:"custInfo,omitempty"` //自定义描述
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
	ruid, _ := strconv.Atoi(c.Request.Header.Get("ruid"))
	// 查询用户sim信息
	info := userSim.Instance().QueryInfo(ruid, req.Sid)
	if info == nil {
		err := fmt.Errorf("not found")
		app.Log.Error(err)
		ret.State, ret.StateInfo = qstate.StateInvalidParameter, err.Error()
		return
	}
	ret.Data = Data{
		Sid:              info.Sid,
		SimNo:            info.SimNo,
		BindTs:           info.BindTs,
		ServiceEndTs:     info.BindTs + 864000, // 绑定后10天
		SimByte:          307200,
		SimAvailableByte: 307200,
	}
	ret.State, ret.StateInfo = State(qstate.StateOk)
}
