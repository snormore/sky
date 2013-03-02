package skyd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
)

const (
	DefaultPort = 8585
)

// A Server is the front end that controls access to tables.
type Server struct {
	httpServer *http.Server
	path       string
	listener   net.Listener
	tables     map[string]*Table
}

// NewServer returns a new Server.
func NewServer(port uint, path string) *Server {
	r := mux.NewRouter()
	s := &Server{
	  httpServer: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r},
	  path:       path,
	  tables:     make(map[string]*Table),
	}

  s.addTableHandlers(r)

	return s
}

// Runs the server.
func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}
	s.listener = listener
	return s.httpServer.Serve(s.listener)
}

// Stops the server.
func (s *Server) Shutdown() error {
	if s.listener != nil {
	  err := s.listener.Close()
	  s.listener = nil
	  return err
	}
	return nil
}

// Processes a request and return the appropriate data format.
func (s *Server) process(w http.ResponseWriter, req *http.Request, f func(params map[string]interface{})(interface{}, error)) {
  // Parses body parameters.
  params := make(map[string]interface{})
  contentType := req.Header.Get("Content-Type")
  if contentType == "application/json" {
	  decoder := json.NewDecoder(req.Body)
  	err := decoder.Decode(&params)
  	if err != nil {
      w.WriteHeader(http.StatusBadRequest)
  		return
  	}
  } else {
    w.WriteHeader(http.StatusNotAcceptable)
    return
  }
  
  // Execute handler.
  ret, err := f(params)
  if err != nil {
    ret = map[string]interface{}{"message":err.Error()}
    w.WriteHeader(500)
  }

  // Encode the return value based on the content type header.
  if ret != nil {
    if contentType == "application/json" {
    	encoder := json.NewEncoder(w)
    	err := encoder.Encode(ret)
    	if err != nil {
    	  w.WriteHeader(http.StatusInternalServerError)
    	  return
    	}
    }
  }

  w.WriteHeader(http.StatusOK)
}

// Generates the path for a table attached to the server.
func (s *Server) GetTablePath(name string) string {
  return fmt.Sprintf("%v/%v", s.path, name)
}

// Retrieves a table that has already been opened.
func (s *Server) GetTable(name string) *Table {
  return s.tables[name]
}
  
// Opens a table and returns a reference to it.
func (s *Server) OpenTable(name string) (*Table, error) {
  // If table already exists then use it.
  table := s.GetTable(name)
  if table != nil {
    return table, nil
  }
  
  // Otherwise open it and save the reference.
  table = NewTable(s.GetTablePath(name))
  err := table.Open()
  if err != nil {
    table.Close()
    return nil, err
  }
  s.tables[name] = table
  
  return table, nil
}
