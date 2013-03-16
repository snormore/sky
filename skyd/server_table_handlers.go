package skyd

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addTableHandlers() {
	s.ApiHandleFunc("/tables", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.createTableHandler(w, req, params)
	}).Methods("POST")
	s.ApiHandleFunc("/tables/{name}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.deleteTableHandler(w, req, params)
	}).Methods("DELETE")
}

// POST /tables
func (s *Server) createTableHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	return s.sync(func() (interface{}, error) {
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
		table = NewTable(tableName, s.TablePath(tableName))
		err = table.Create()
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}

// DELETE /tables/:name
func (s *Server) deleteTableHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	tableName := vars["name"]

	return s.sync(func() (interface{}, error) {
		// Return an error if the table doesn't exist.
		table := s.GetTable(tableName)
		if table == nil {
			table = NewTable(tableName, s.TablePath(tableName))
		}
		if !table.Exists() {
			return nil, errors.New("Table does not exist.")
		}

		// Otherwise delete it.
		return nil, table.Delete()
	})
}
