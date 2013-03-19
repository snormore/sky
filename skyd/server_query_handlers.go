package skyd

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addQueryHandlers() {
	s.ApiHandleFunc("/tables/{name}/query", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryHandler(w, req, params)
	}).Methods("POST")
	s.ApiHandleFunc("/tables/{name}/query/codegen", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryCodegenHandler(w, req, params)
	}).Methods("POST")
}

// POST /tables/:name/query
func (s *Server) queryHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.RunQuery(vars["name"], params)
}

// POST /tables/:name/query/codegen
func (s *Server) queryCodegenHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Retrieve table and codegen query.
	var source string
	_, err := s.sync(func() (interface{}, error) {
		// Return an error if the table already exists.
		table, err := s.OpenTable(vars["name"])
		if err != nil {
			return nil, err
		}

		// Deserialize the query.
		query := NewQuery(table)
		err = query.Deserialize(params)
		if err != nil {
			return nil, err
		}

		// Generate the query source code.
		source, err = query.Codegen()
		fmt.Println(source)
		return source, err
	})

	if err == nil {
		return map[string]interface{}{"source": source}, nil
	}

	return nil, err
}
