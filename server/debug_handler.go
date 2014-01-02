package server

import (
	"expvar"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

func installDebugHandler(s *Server) {
	s.Router.HandleFunc("/debug/pprof", pprof.Index)
	s.Router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.Router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.Router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.Router.HandleFunc("/debug/pprof/{name}", pprof.Index)

	s.Router.HandleFunc("/debug/vars", getVarsHandler)
	s.Router.HandleFunc("/debug/vars/{name}", getVarsHandler)
}

func getVarsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	expvar.Do(func(kv expvar.KeyValue) {
		if vars["name"] == "" || vars["name"] == kv.Key {
			if !first {
				fmt.Fprintf(w, ",\n")
			}
			first = false
			fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
		}
	})
	fmt.Fprintf(w, "\n}\n")
}
