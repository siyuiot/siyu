package httpserver

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

// 初始化连接
func initClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.0.247:30279",
		Password: "", // no password set
		DB:       4,  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

type r struct {
	Did        int  `json:"did"`         // 车子bike_id
	Dts        int  `json:"dts"`         // 时间戳
	Ride       bool `json:"ride"`        // 车辆是否处于骑行模式
	Acc        bool `json:"acc"`         // 是否通电
	Defence    int  `json:"defence"`     // 设防状态 1设防 2静防？3不设防
	SeaterLock bool `json:"seater_lock"` // 座椅锁，1打开，2锁定
	Vol        int  `json:"vol"`         // 电池电压
	TotalMiles int  `json:"total_miles"` // 总距离 单位km
}


func TestQueryInfo(t *testing.T) {
	_ = initClient()
	info, err := rdb.HGetAll("br:10000007").Result()
	if err != nil && len(info) < 1 {
		return
	}
	acc := false
	accInfo := info["acc"]
	if accInfo == "1" {
		acc = true
	}
	vol, _ := strconv.Atoi(info["vol"])

	total := 0
	miles := info["miles"]
	if len(miles) > 0 {
		m := make(map[string]int)
		_ = json.Unmarshal([]byte(miles), &m)
		total = m["total"]
	}
	fmt.Println(total)
	fmt.Println(reflect.TypeOf(total))

	res := r{
		Did:        10000007,
		Dts:        0,
		Ride:       false,
		Acc:        acc,
		Defence:    0,
		SeaterLock: false,
		Vol:        vol,
		TotalMiles: total,
	}
	fmt.Println(res)
}
