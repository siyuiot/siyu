package qgin

import (
	"bytes"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/siyuiot/siyu/pkg/qlog"
)

type resWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r resWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

type logOption struct {
	MaxResLen int // 最大返回值长度,超出不打印 log
}

var defaultOpt = logOption{
	MaxResLen: 4096,
}

type optionFunc func(opts *logOption)

func WithMaxResLen(len int) optionFunc {
	return func(opts *logOption) {
		opts.MaxResLen = len
	}
}

func LogReq(l *qlog.Entry, opts ...optionFunc) gin.HandlerFunc {
	opt := defaultOpt
	for _, apply := range opts {
		apply(&opt)
	}
	return func(c *gin.Context) {
		start := time.Now()
		w := &resWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		contentType := c.Request.Header.Get("Content-Type")
		_log := l.WithField(MiddleReqId, GetRequestId(c))
		uid := GetUid(c)
		if uid > 0 {
			_log = l.WithField("uid", GetUid(c))
		}
		if !strings.Contains(contentType, "multipart/form-data") {
			_log.Infof("url=%s req=%s ip=%s", c.Request.URL, bodyBytes, RemoteIp(c.Request))
		} else {
			_log.Infof("url=%s", c.Request.URL)
		}
		_ = c.Request.Body.Close() //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		l.WithField("latency", latency)
		if len(w.body.Bytes()) < opt.MaxResLen {
			_log.Debugf("url=%s res=%s", c.Request.URL, w.body.Bytes())
		} else {
			_log.Debugf("url=%s res too long", c.Request.URL)
		}
	}
}

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}
