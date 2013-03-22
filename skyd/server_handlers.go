package skyd

import (
	"net/http"
)

func (s *Server) addHandlers() {
	s.ApiHandleFunc("/ping", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.pingHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/rawping", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.rawPingHandler(w, req, params)
	}).Methods("GET")
}

// GET /ping
func (s *Server) pingHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	return s.sync(func() (interface{}, error) {
		return map[string]interface{}{"message": "ok"}, nil
	})
}

// GET /rawping
func (s *Server) rawPingHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}
