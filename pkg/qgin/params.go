package qgin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/modules/qstate"
)

const (
	MiddleSgagent      = "Sgagent"
	MiddleAppKey       = "app"
	MiddleUidKey       = "ruserid"   // uid
	MiddleBattUidKey   = "battuid"   // 换电uid
	MiddleBattAdminKey = "battadmin" // 换电admin uid
	MiddleAdminUiKey   = "adminuid"  // 管理员uid
	MiddleAdminKey     = "daccount"  // 管理员用户名
	MiddleReqId        = "requestId"
	MiddleAreaCode     = "areacode"
	MiddleLang         = "lang"
	MiddleBattAppkey   = "appkey"
	MiddleAppverkey    = "appver"   // app版本
	MiddlePlatformkey  = "platform" // app平台
	MiddlePkgNamekey   = "pkgname"  // pkgname
	MiddleSysverkey    = "sysver"   // SysVersion
	MiddleDidkey       = "appdid"   // 设备号
)

type SingleId struct {
	qstate.CommonRequest
	Id int `json:"id" binding:"required"`
}

type PageParams struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func GetUidAndReqId(c *gin.Context) map[string]string {
	headers := map[string]string{
		MiddleUidKey: getuid(c),
		MiddleReqId:  GetRequestId(c)}
	return headers
}

func GetRequestId(c *gin.Context) string {
	return c.Request.Header.Get(MiddleReqId)
}

func GetUid(c *gin.Context) int {
	ruserid, _ := strconv.Atoi(getuid(c)) //网关层保证ruserid有效
	return ruserid
}

func GetBattUid(c *gin.Context) int {
	uid := c.Request.Header.Get(MiddleBattUidKey)
	ruserid, _ := strconv.Atoi(uid)
	return ruserid
}

func GetAreaCode(c *gin.Context) string {
	return c.Request.Header.Get(MiddleAreaCode)
}

func GetLang(c *gin.Context) string {
	return c.Request.Header.Get(MiddleLang)
}

func GetBattAdminUid(c *gin.Context) int {
	uid := c.Request.Header.Get(MiddleBattAdminKey)
	ruserid, _ := strconv.Atoi(uid)
	return ruserid
}

func GetApp(c *gin.Context) string {
	return c.Request.Header.Get(MiddleAppKey)
}

func GetPlatform(c *gin.Context) string {
	return c.Request.Header.Get(MiddlePlatformkey)
}

func GetPkgName(c *gin.Context) string {
	return c.Request.Header.Get(MiddlePkgNamekey)
}

func GetAppVer(c *gin.Context) string {
	return c.Request.Header.Get(MiddleAppverkey)
}

func GetSysver(c *gin.Context) string {
	return c.Request.Header.Get(MiddleSysverkey)
}

func GetAppDid(c *gin.Context) string {
	return c.Request.Header.Get(MiddleDidkey)
}

func GetAdminUid(c *gin.Context) int {
	uid := c.Request.Header.Get(MiddleAdminUiKey)
	adminUid, _ := strconv.Atoi(uid)
	return adminUid
}

func GetAdminAccount(c *gin.Context) string {
	return c.Request.Header.Get(MiddleAdminKey)
}

func getuid(c *gin.Context) string {
	return c.Request.Header.Get(MiddleUidKey)
}

func QueryInt(c *gin.Context, key string, defaultVal int) int {
	val, ok := c.GetQuery(key)
	if ok {
		res, _ := strconv.Atoi(val)
		return res
	}
	return defaultVal
}
func QueryStr(c *gin.Context, key string, defaultVal string) string {
	val, ok := c.GetQuery(key)
	if ok {
		return val
	}
	return defaultVal
}
func ParamInt(c *gin.Context, key string, defaultVal int) int {
	val := c.Param(key)
	if val == "" {
		res, _ := strconv.Atoi(val)
		return res
	}
	return defaultVal
}
func ParamStr(c *gin.Context, key string, defaultVal string) string {
	val := c.Param(key)
	if val == "" {
		return val
	}
	return defaultVal
}
