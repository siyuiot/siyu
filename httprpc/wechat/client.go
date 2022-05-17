package wechat

import "github.com/siyuiot/siyu/pkg/qhttp"

const (
	code2Session = "/wechat/code2Session"
)

var (
	defaultAddr = "http://wechat.sg-base:8080"
)

type Client struct {
	addr      string
	requestId string
}

type ClientOption struct {
	Addr      string
	RequestId string
}

func New(in ClientOption) *Client {
	if in.Addr == "" {
		in.Addr = defaultAddr
	}
	return &Client{addr: in.Addr, requestId: in.RequestId}
}

func (c *Client) Code2Session(in InCode2Session) (Code2SessionData, error) {
	var out = &OutCode2Session{}
	headers := map[string]string{"requestId": c.requestId}
	_, err := qhttp.PostJSON(c.addr+code2Session, headers, in, &out)
	return out.Data, err
}
