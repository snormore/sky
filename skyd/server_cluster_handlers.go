package skyd

import (
	"net/http"
)

func (s *Server) addClusterHandlers() {
	s.ApiHandleFunc("/cluster", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getClusterHandler(w, req, params)
	}).Methods("GET")
}

// GET /cluster
func (s *Server) getClusterHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	return s.cluster.serialize(), nil
}

// POST /cluster/groups
func (s *Server) createClusterNodeGroupHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	command := &CreateNodeGroupCommand{NodeGroupId: NewNodeGroupId()}
	return nil, s.Do(command)
}
