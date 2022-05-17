package main

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/siyuiot/siyu/pkg/qlog"
)

func main() {
	time.Sleep(time.Second)
	log.Debug("This is a DEBUG message")
	log.Info("This is a INFO message")
	log.Warn("This is a WARN message")
	log.Error("This is a ERROR message")
	//log.Fatal("This is a FATAL message")
	//log.Panic("This is a PANIC message")
	log.WithField("foo", "bar").Warn("This is a WithField WARN message")
	log.WithField("hello", "world").WithField("goobye", "moon").Info("This is a WithFields INFO message")
	log.WithError(errors.New("An error")).Warn("with error warning")

	entry := log.WithField("test", "1")
	entry.Debug("This is a DEBUG message from entry")
	entry.Info("This is a INFO message from entry")

	// try to change the config now
	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		for i := 0; i < 100; i++ {
			ok(entry)
			time.Sleep(time.Second)
		}
		cancel()
	}()

	<-ctx.Done()
}

func ok(entry *logrus.Entry) {
	entry.WithField("app", "log").WithField("event", "1111").Warn("hello ", 1)
}
