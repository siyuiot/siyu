package qstate

type State = int

// 按业务模块拆分状态码
// 此处只定义业务必须状态码,服务端故障情况太多,都返出去无意义
// 不定义http300,http400,http500
// 不定义rpc错误
// 不定义db层错误
// 特殊情况除外,如04数据库无记录
// golbal 00~100
// gateway 1000~1100
// user 1100~1200
// bikesvc 1200~1300
// order   1300~1500
// ota升级  1500~1600
// admin  1600~1800
// >100个模块要拆分了
const (
	//golbal 00~100
	StateOk                  = 0 //"00" // 成功
	StateInvalidParameter    = 1 //"01" // 参数错误
	StateFailed              = 2 //"02" // 失败
	StateCrossDomainPrelight = 3 //"03" // 跨域预检
	StateNoRecord            = 4 //"04" // 无记录
	StateDuplicate           = 5 //"05" // 重复创建
	StateNotAuth             = 6 //"06" // 未认证
)

var defaultStateInfo = map[State]string{
	StateOk:               "ok",
	StateInvalidParameter: "invalid params",
	StateFailed:           "failed",
}

func StateStr(s State) string {
	return defaultStateInfo[s]
}

type CommonRequest struct {
	Ts  int64
	Seq int64
}

type CommonResponse struct {
	State     int         `json:"state"`               //状态值
	StateInfo string      `json:"state_info"`          //状态描述
	CustInfo  string      `json:"cust_info,omitempty"` //自定义描述
	Seq       int         `json:"seq"`
	Data      interface{} `json:"data"`
}
