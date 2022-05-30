package userSim

type Info struct {
	Uid              int
	Sid              int
	SimProvider      string
	SimNo            string
	Iccid            string
	BindTs           int64
	ServiceEndTs     int64 // 服务结束时间戳
	ServiceDuration  int   // 服务时长
	SimByte          int   // 卡流量
	SimAvailableByte int   // 可用流量
	Remark           string
	Created          int64
	Updated          int64
}
