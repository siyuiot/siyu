package qhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	QHTTPTransport http.RoundTripper = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   5,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	QHTTPClient = &http.Client{
		Transport: QHTTPTransport,
	}
)

// 保持http链接复用，只是resp.Body.Close()不够
// 	1. ioutil.ReadAll(resp.Body) && resp.Body.Close()
// 	2. io.Copy(ioutil.Discard,resp.Body) && resp.Body.Close()
func Get(url string) (resp *http.Response, err error) {
	return QHTTPClient.Get(url)
}

// 保持http链接复用，只是resp.Body.Close()不够
// 	1. ioutil.ReadAll(resp.Body) && resp.Body.Close()
// 	2. io.Copy(ioutil.Discard,resp.Body) && resp.Body.Close()
func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return QHTTPClient.Post(url, contentType, body)
}

// 保持http链接复用，只是resp.Body.Close()不够
// 	1. ioutil.ReadAll(resp.Body) && resp.Body.Close()
// 	2. io.Copy(ioutil.Discard,resp.Body) && resp.Body.Close()
func PostWithHeader(url string, headers http.Header, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v[0])
		}
	}
	rsp, err := QHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return rsp, err
}

// 保持http链接复用，只是resp.Body.Close()不够
// 	1. ioutil.ReadAll(resp.Body) && resp.Body.Close()
// 	2. io.Copy(ioutil.Discard,resp.Body) && resp.Body.Close()
func PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return QHTTPClient.PostForm(url, data)
}

// 保持http链接复用，只是resp.Body.Close()不够
// 	1. ioutil.ReadAll(resp.Body) && resp.Body.Close()
// 	2. io.Copy(ioutil.Discard,resp.Body) && resp.Body.Close()
func PostFormWithHeader(url string, headers map[string]string, data url.Values) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	return QHTTPClient.Do(req)
}

func Head(url string) (resp *http.Response, err error) {
	return QHTTPClient.Head(url)
}

func PostJSON(url string, headers map[string]string, body interface{}, result interface{}) (resp *http.Response, err error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	rsp, err := QHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if result != nil { //if result != nil, try Unmarshal the body
		defer rsp.Body.Close()
		rspBody, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(rspBody, result)
		if err != nil {
			return rsp, err
		}
	}
	return rsp, err
}

func GetJSON(url string, result interface{}) (data []byte, err error) {
	resp, err := QHTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	rspBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rspBody, result)
	if err != nil {
		return nil, err
	}
	return rspBody, nil
}

func Download(data []byte, w http.ResponseWriter, r *http.Request, filename string) {
	w.Header().Set("Content-Type", "application/octet-stream;charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))
	w.Write(data)
}
