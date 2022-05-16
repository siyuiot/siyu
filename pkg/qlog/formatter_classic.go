package qlog

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

//copy from old logrus code
func prefixFieldClashes(data logrus.Fields, hasCaller bool) {
	if t, ok := data["time"]; ok {
		data["fields.time"] = t
	}

	if m, ok := data["msg"]; ok {
		data["fields.msg"] = m
	}

	if l, ok := data["level"]; ok {
		data["fields.level"] = l
	}

	if hasCaller {
		if fc, ok := data["func"]; ok {
			data["fields.func"] = fc
		}

		if f, ok := data["file"]; ok {
			data["fields.file"] = f
		}
	}
}

// ShortLevel is a simple char to indicate log level
type ShortLevel uint32

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level ShortLevel) String() string {
	_level := logrus.Level(level)
	switch _level {
	case logrus.TraceLevel:
		return "T"
	case logrus.DebugLevel:
		return "D"
	case logrus.InfoLevel:
		return "I"
	case logrus.WarnLevel:
		return "W"
	case logrus.ErrorLevel:
		return "E"
	case logrus.FatalLevel:
		return "F"
	case logrus.PanicLevel:
		return "P"
	}

	return "X"
}

// ClassicFormatter formats logs into parsable json
type ClassicFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	TruncateCallerPath bool

	CallerPathStrip bool
	SortField       bool
	DisableSorting  bool
	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)
}

func stripPathAuto(path string) string {
	//filters := []string{"/github.com/", "/app/", "/pkg/", "/grpc/", "/modules/", "/httprpc/"}
	parts := strings.Split(path, "/")
	plen := len(parts)
	if plen > 2 {
		return strings.Join(parts[plen-2:], "/")
	}
	return path
}

// Format renders a single log entry
func (f *ClassicFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	prefixFieldClashes(entry.Data, entry.HasCaller())
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//timestamp
	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = longTimeStamp
	}
	if !f.DisableTimestamp {
		b.WriteString(entry.Time.Format(timestampFormat))
	}
	//timestamp level
	f.appendValueOnly(b, fmt.Sprintf("%s", ShortLevel(entry.Level).String()))
	//timestamp level app shortFile/fileName:line
	server := "app"
	if entry.HasCaller() { //reportcaller is enabled
		var app string = "-"
		if v, ok := entry.Data[server]; ok {
			app = v.(string)
		}
		if f.TruncateCallerPath {
			if !f.CallerPathStrip {
				f.appendValueOnly(b, fmt.Sprintf("%s %s:%d", app, filepath.Base(entry.Caller.File), entry.Caller.Line))
			} else {
				f.appendValueOnly(b, fmt.Sprintf("%s %s:%d", app, stripPathAuto(entry.Caller.File), entry.Caller.Line))
			}
		} else {
			f.appendValueOnly(b, fmt.Sprintf("%s %s:%d", app, entry.Caller.File, entry.Caller.Line))
		}
	}
	//timestamp level app shortFile/fileName:line message
	if len(entry.Message) > 0 {
		f.appendValueOnly(b, entry.Message)
	}
	//timestamp level app shortFile/fileName:line message filed...
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	filedKeys := make([]string, 0, len(entry.Data))
	if !f.DisableSorting { // field 排序
		if f.SortingFunc == nil {
			sort.Strings(keys)
			filedKeys = append(filedKeys, keys...)
		} else {
			filedKeys = append(filedKeys, keys...)
			f.SortingFunc(filedKeys)
		}
	} else {
		filedKeys = append(filedKeys, keys...)
	}
	for _, k := range filedKeys {
		if k == server {
			continue
		}
		f.appendKeyValue(b, k, entry.Data[k])
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *ClassicFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *ClassicFormatter) appendValueOnly(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	f.appendValue(b, value)
}

func (f *ClassicFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	b.WriteString(fmt.Sprintf("%s", stringVal))
}

var _InitClassicFormatter = func() interface{} {
	registeFormatter("classic", reflect.TypeOf(ClassicFormatter{}))
	return nil
}()
