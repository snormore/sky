package skyd

import (
	"encoding/json"
	"github.com/benbjohnson/go-raft"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addRaftHandlers() {
	s.ApiHandleFunc("/raft/run/{name}", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.doRaftCommandHandler(w, req, params)
	}).Methods("POST")
}

// POST /raft
func (s *Server) doRaftCommandHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	command, err := raft.NewCommand(vars["name"])
	if err != nil {
		return nil, err
	}
	json.NewDecoder(req.Body).Decode(&command)

	return nil, s.Do(command)
}
