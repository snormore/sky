package skyd

import (
	"net/http"
)

func (s *Server) addHandlers() {
	s.ApiHandleFunc("/ping", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.pingHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.indexHandler(w, req, params)
	}).Methods("GET")
}

// GET /ping
func (s *Server) pingHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	return map[string]interface{}{"message": "ok"}, nil
}

// GET /
func (s *Server) indexHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	return map[string]interface{}{"sky": "welcome", "version": Version}, nil
}
