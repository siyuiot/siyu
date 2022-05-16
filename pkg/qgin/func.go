package qgin

import (
	"fmt"

	"github.com/spf13/viper"
)

type runMode = string

const (
	modeDev   runMode = "dev"
	modeTest  runMode = "test"
	modeDebug runMode = "debug"
	modeProd  runMode = "prod"
)
const RedisDevUrl = "192.168.0.247:30279"
const RedisDevAuth = ""
const DbUserUrl = "postgres://postgres:iLoveShark@192.168.0.247:32432/bsmaster?sslmode=disable&fallback_application_name=user"

const DbOrderUrl = "postgres://postgres:iLoveShark@192.168.0.247:32432/bsorder?sslmode=disable&fallback_application_name=order"
const DbLogUrl = "postgres://postgres:iLoveShark@192.168.0.247:32432/bslog?sslmode=disable&fallback_application_name=operationlog"

var TsdbBikeRideUrl = GetTsdbUrl("bsbikeride", "")

const dburl = "postgres://postgres:iLoveShark@192.168.0.247:32432/%s?sslmode=disable&fallback_application_name=%s"
const tsdb = "postgres://postgres:iLoveShark@192.168.0.247:32433/%s?sslmode=disable&fallback_application_name=%s"

func GetDbUrl(db, app string) string {
	return fmt.Sprintf(dburl, db, app)
}

func GetTsdbUrl(db, app string) string {
	return fmt.Sprintf(tsdb, db, app)
}

func IsDev() bool {
	return getRunMode() == modeDev
}

func getRunMode() string {
	return viper.GetString("runmode")
}
func IsTest() bool {
	return getRunMode() == modeTest
}
func IsDebug() bool {
	return getRunMode() == modeDebug
}
func IsProd() bool {
	return getRunMode() == modeProd
}
