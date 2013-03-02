package skyd

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addTableHandlers(r *mux.Router) {
	r.HandleFunc("/tables", func(w http.ResponseWriter, req *http.Request) { s.createTableHandler(w, req) }).Methods("POST")
}

// POST /tables
func (s *Server) createTableHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    // Retrieve table parameters.
    tableName, ok := params["name"].(string)
    if !ok {
      return nil, errors.New("Table name required.")
    }
    
    // Return an error if the table already exists.
    table, err := s.OpenTable(tableName)
    if table != nil {
      return nil, errors.New("Table already exists.")
    }
    
    // Otherwise create it.
    table = NewTable(s.GetTablePath(tableName))
    err = table.Create()
    if err != nil {
      return nil, err
    }
    
  	return nil, nil
  })
}
