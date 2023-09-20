package pprof

import (
	"net/http"
	netProf "net/http/pprof"
	"strings"
)

// Route registers pprof routes based on the basepath
func Route(basePath string) {
	http.HandleFunc(basePath+"/debug/pprof/", index(basePath))
	http.HandleFunc(basePath+"/debug/pprof/cmdline", netProf.Cmdline)
	http.HandleFunc(basePath+"/debug/pprof/profile", netProf.Profile)
	http.HandleFunc(basePath+"/debug/pprof/symbol", netProf.Symbol)
	http.HandleFunc(basePath+"/debug/pprof/trace", netProf.Trace)
}

func index(basePath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, basePath)
		netProf.Index(w, r)
	}
}
