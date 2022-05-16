package qhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// go test -v -test.run TestHttpKeepAlive
// while true;do echo $(date --rfc-3339=seconds);netstat -apn|grep "210.22.78.6:80";sleep 1s;done
func TestHttpKeepAlive(t *testing.T) {
	url := "http://api.d.blueshark.com"
	// data := map[string]string{"foo", "bar"}
	data := []byte("foobar")
	for i := 0; i < 10; i++ {
		resp, err := PostWithHeader(url, nil, bytes.NewReader(data))
		if err != nil {
			fmt.Printf("call url(%s) error(%v)\n", url, err)
			return
		}
		io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("call url(%s) response error(%d)\n", url, resp.StatusCode)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("call url(%s) response error(%d)\n", url, resp.StatusCode)
			return
		}
		fmt.Println(time.Now(), resp.StatusCode, string(body))
		time.Sleep(1 * time.Second)
	}
}
