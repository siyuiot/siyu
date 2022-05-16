package qdebugserver

import (
	"encoding/json"
	"expvar"
	"io"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

// RegisterHTTPMux register HTTP mux
func RegisterHTTPMux(mux *http.ServeMux, prefixOptions ...string) *http.ServeMux {
	prefix := getPrefix(prefixOptions...)

	mux.HandleFunc(prefix+"/", debugIndex)
	mux.HandleFunc(prefix+"/pprof/", pprof.Index)
	mux.HandleFunc(prefix+"/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc(prefix+"/pprof/profile", pprof.Profile)
	mux.HandleFunc(prefix+"/pprof/symbol", pprof.Symbol)
	mux.HandleFunc(prefix+"/pprof/trace", pprof.Trace)

	mux.Handle(prefix+"/vars", expvar.Handler())

	mux.HandleFunc(prefix+"/healthz", healthHandler)
	mux.HandleFunc(prefix+"/readyz", readyzHandler)
	mux.HandleFunc(prefix+"/version", versionHandler)

	return mux
}

// copy from https://github.com/gin-contrib/pprof/
const (
	// DefaultPrefix url prefix of pprof
	DefaultPrefix = "/debug"
)

func getPrefix(prefixOptions ...string) string {
	prefix := DefaultPrefix
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

// RegisterGin the standard HandlerFuncs from the net/http/pprof package with
// the provided gin.Engine. prefixOptions is a optional. If not prefixOptions,
// the default path prefix is used, otherwise first prefixOptions will be path prefix.
func RegisterGin(r *gin.Engine, prefixOptions ...string) *gin.RouterGroup {
	prefix := getPrefix(prefixOptions...)

	debugGroup := r.Group(prefix)
	{
		debugGroup.GET("/", pprofHandler(debugIndex))

		prefixPprof := debugGroup.Group("/pprof")
		{
			prefixPprof.GET("/", pprofHandler(pprof.Index))
			prefixPprof.GET("/cmdline", pprofHandler(pprof.Cmdline))
			prefixPprof.GET("/profile", pprofHandler(pprof.Profile))
			prefixPprof.POST("/symbol", pprofHandler(pprof.Symbol))
			prefixPprof.GET("/symbol", pprofHandler(pprof.Symbol))
			prefixPprof.GET("/trace", pprofHandler(pprof.Trace))
			prefixPprof.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
			prefixPprof.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
			prefixPprof.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
			prefixPprof.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
			prefixPprof.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
			prefixPprof.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
		}
		debugGroup.GET("/vars", pprofHandler(expvar.Handler().ServeHTTP))

		debugGroup.GET("/healthz", pprofHandler(healthHandler))
		debugGroup.GET("/readyz", pprofHandler(readyzHandler))
		debugGroup.GET("/version", pprofHandler(versionHandler))
	}
	return debugGroup
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

var (
	userReadyzHandler http.HandlerFunc
	versions          map[string]string
)

// SetUserReadyzHandler set user specified readyz handler
func SetUserReadyzHandler(handleFunc http.HandlerFunc) {
	userReadyzHandler = handleFunc
}

// SetVersionInfo set app version
func SetVersionInfo(ver map[string]string) {
	versions = ver
}

func readyzHandler(w http.ResponseWriter, r *http.Request) {
	if userReadyzHandler != nil {
		userReadyzHandler(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	s, err := json.Marshal(versions)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "marshal error:")
		io.WriteString(w, err.Error())
		return
	}
	w.Write(s)
}
