package userSim

type Info struct {
	Uid              int
	Sid              int
	SimNo            string
	BindTs           int64
	ServiceEndTs     int64 // 服务结束时间戳
	SimByte          int   // 卡流量
	SimAvailableByte int   // 可用流量
	Remark           string
	Created          int64
	Updated          int64
}
