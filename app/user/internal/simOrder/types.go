package simOrder

type orderStatus = int // 订单状态
const (
	PaymentSubmit  orderStatus = 1
	PaymentSuccess orderStatus = 2
)

type Info struct {
	Oid         int
	Uid         int
	Sid         int
	Typ         int
	No          string
	Name        string
	SkuId       int
	Status      int
	AmountPrice int64
	DuePrice    int64
	PayPrice    int64
	PayChannel  string
	Remark      string
	Created     int64
	Updated     int64
}
