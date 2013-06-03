package skyd

import (
	"errors"
	"fmt"
	"github.com/benbjohnson/go-raft"
	"net/http"
)

func (s *Server) addClusterHandlers() {
	s.ApiHandleFunc("/cluster", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.getClusterHandler(w, req, params)
	}).Methods("GET")

	s.ApiHandleFunc("/cluster/commands", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.doClusterCommandHandler(w, req, params)
	}).Methods("POST")

	s.ApiHandleFunc("/cluster/nodes", &CreateNodeCommand{}, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.createClusterNodeHandler(w, req, params)
	}).Methods("POST")
}

// GET /cluster
func (s *Server) getClusterHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	return s.cluster.serialize(), nil
}

// POST /cluster/commands
func (s *Server) doClusterCommandHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	command := params.(raft.Command)
	return nil, s.raftServer.Do(command)
}

// POST /cluster/nodes
func (s *Server) createClusterNodeHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	command := params.(*CreateNodeCommand)
	
	// Generate a node id if one is not passed in.
	if command.NodeId == "" {
		command.NodeId = NewNodeId()
	}

	// Retrieve a group to add to if one is not specified.
	if command.NodeGroupId == "" {
		group := s.cluster.GetAvailableNodeGroup()
		if group == nil {
			return nil, errors.New("No groups available")
		}
	}

	// Require host & port.
	if command.Host == "" {
		return nil, fmt.Errorf("Invalid host: %s", command.Host)
	}
	if command.Port == 0 {
		return nil, fmt.Errorf("Invalid port: %d", command.Port)
	}

	return nil, s.Do(command)
}
