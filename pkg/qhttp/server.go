package qhttp

// import (
// 	"context"
// 	"net/http"
// 	"time"

// 	"gogs.sharkgulf.cn/sg/bs/pkg/qlog"
// )

// // RunServer run a http server with gracefully shutdown
// func RunServer(l *qlog.Entry, ctx context.Context, addr string, handler http.Handler) (err error) {
// 	_log := l
// 	_log.Debugf("Listening and serving HTTP on %s\n", addr)
// 	srv := &http.Server{
// 		Addr:    addr,
// 		Handler: handler,
// 	}
// 	srvErrChan := make(chan error)
// 	go func() {
// 		// service connections
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			_log.WithError(err).Error("HTTP server error")
// 			srvErrChan <- err
// 		}
// 	}()
// 	select {
// 	case err = <-srvErrChan:
// 		return err
// 	case <-ctx.Done():
// 		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		if err = srv.Shutdown(shutDownCtx); err != nil {
// 			_log.WithError(err).Error("Server Shutdown:", err)
// 			return err
// 		}
// 	}
// 	_log.Debug("Server gracefully exit")
// 	return nil
// }
