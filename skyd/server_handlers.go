package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addHandlers(r *mux.Router) {
	r.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) { s.pingHandler(w, req) }).Methods("GET")
	r.HandleFunc("/rawping", func(w http.ResponseWriter, req *http.Request) { s.rawPingHandler(w, req) }).Methods("GET")
}

// GET /ping
func (s *Server) pingHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    return map[string]interface{}{"message":"ok"}, nil
  })
}

// GET /rawping
func (s *Server) rawPingHandler(w http.ResponseWriter, req *http.Request) {
}
