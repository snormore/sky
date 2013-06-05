package skyd

import (
	"errors"
	"github.com/benbjohnson/go-raft"
	"net/http"
)

func (s *Server) addClusterHandlers() {
	s.ApiHandleFunc("/cluster", nil, s.getClusterHandler).Methods("GET")
	s.ApiHandleFunc("/cluster/commands", nil, s.clusterExecuteCommandHandler).Methods("POST")
	s.ApiHandleFunc("/cluster/append", &raft.AppendEntriesRequest{}, s.clusterAppendEntriesHandler).Methods("POST")
	s.ApiHandleFunc("/cluster/nodes", &CreateNodeCommand{}, s.clusterCreateNodeHandler).Methods("POST")
	s.ApiHandleFunc("/cluster/nodes", &RemoveNodeCommand{}, s.clusterRemoveNodeHandler).Methods("DELETE")
}

// GET /cluster
func (s *Server) getClusterHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	return s.cluster.serialize(), nil
}

// POST /cluster/commands
func (s *Server) clusterExecuteCommandHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	command := params.(raft.Command)
	return nil, s.ExecuteClusterCommand(command)
}

// POST /cluster/append
func (s *Server) clusterAppendEntriesHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	r := params.(*raft.AppendEntriesRequest)

	// If the log has not been appended to (except for server init) then
	// truncate it and allow entries. This occurs in the case of a server join.
	if s.ClusterRaftMemberCount() == 1 {
		if err := s.Reset(); err != nil {
			return nil, err
		}
	}

	// Retrieve the Raft server.
	s.mutex.Lock()
	raftServer := s.clusterRaftServer
	s.mutex.Unlock()
	if raftServer == nil {
		return nil, errors.New("skyd: Raft server unavailable")
	}

	resp, err := s.clusterRaftServer.AppendEntries(r)
	if err != nil {
		warn("[/append] %v", err)
	}
	return resp, nil
}

// POST /cluster/nodes
func (s *Server) clusterCreateNodeHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	command := params.(*CreateNodeCommand)

	// Retrieve a group to add to if one is not specified.
	if command.NodeGroupId == "" {
		group := s.cluster.GetAvailableNodeGroup()
		if group == nil {
			return nil, errors.New("No groups available")
		}
		command.NodeGroupId = group.id
	}

	return nil, s.ExecuteClusterCommand(command)
}

// DELETE /cluster/nodes
func (s *Server) clusterRemoveNodeHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	command := params.(*RemoveNodeCommand)
	return nil, s.ExecuteClusterCommand(command)
}
