package skyd

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addTableHandlers() {
	s.ApiHandleFunc("/tables", nil, s.getTablesHandler).Methods("GET")
	s.ApiHandleFunc("/tables/{name}", nil, s.getTableHandler).Methods("GET")
	s.ApiHandleFunc("/tables", nil, s.createTableHandler).Methods("POST")
	s.ApiHandleFunc("/tables/{name}", nil, s.deleteTableHandler).Methods("DELETE")
}

// GET /tables
func (s *Server) getTablesHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	return s.GetAllTables()
}

// GET /tables/:name
func (s *Server) getTableHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.OpenTable(vars["name"])
}

// POST /tables
func (s *Server) createTableHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	p := params.(map[string]interface{})

	// Retrieve table parameters.
	tableName, ok := p["name"].(string)
	if !ok {
		return nil, errors.New("Table name required.")
	}

	// Return an error if the table already exists.
	table, err := s.OpenTable(tableName)
	if table != nil {
		return nil, errors.New("Table already exists.")
	}

	// Otherwise create it.
	table = NewTable(tableName, s.TablePath(tableName))
	err = table.Create()
	if err != nil {
		return nil, err
	}

	return table, nil
}

// DELETE /tables/:name
func (s *Server) deleteTableHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	tableName := vars["name"]

	return nil, s.DeleteTable(tableName)
}
