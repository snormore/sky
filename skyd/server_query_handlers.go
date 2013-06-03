package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addQueryHandlers() {
	s.ApiHandleFunc("/tables/{name}/stats", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.statsHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/query", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.queryHandler(w, req, params)
	}).Methods("POST")
	s.ApiHandleFunc("/tables/{name}/query/codegen", nil, func(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
		return s.queryCodegenHandler(w, req, params)
	}).Methods("POST")
}

// GET /tables/:name/stats
func (s *Server) statsHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Run a simple count query.
	query := NewQuery(table, s.factors)
	selection := NewQuerySelection(query)
	selection.Fields = append(selection.Fields, NewQuerySelectionField("count", "count()"))
	query.Prefix = req.FormValue("prefix")
	query.Steps = append(query.Steps, selection)

	return s.RunQuery(table, query)
}

// POST /tables/:name/query
func (s *Server) queryHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	p := params.(map[string]interface{})
	
	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Deserialize the query.
	query := NewQuery(table, s.factors)
	err = query.Deserialize(p)
	if err != nil {
		return nil, err
	}

	return s.RunQuery(table, query)
}

// POST /tables/:name/query/codegen
func (s *Server) queryCodegenHandler(w http.ResponseWriter, req *http.Request, params interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	p := params.(map[string]interface{})

	// Retrieve table and codegen query.
	var source string
	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Deserialize the query.
	query := NewQuery(table, s.factors)
	err = query.Deserialize(p)
	if err != nil {
		return nil, err
	}

	// Generate the query source code.
	source, err = query.Codegen()
	//fmt.Println(source)

	return source, &TextPlainContentTypeError{}
}
