package app

import "github.com/siyuiot/siyu/pkg/qlog"

const Name = "user"

var Log = qlog.WithField("app", Name)
