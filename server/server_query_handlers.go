package server

import (
	"github.com/gorilla/mux"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query"
	"net/http"
)

func (s *Server) addQueryHandlers() {
	s.ApiHandleFunc("/tables/{name}/stats", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.statsHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/query", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryHandler(w, req, params)
	}).Methods("POST")
	s.ApiHandleFunc("/tables/{name}/query/codegen", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryCodegenHandler(w, req, params)
	}).Methods("POST")
}

// GET /tables/:name/stats
func (s *Server) statsHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	// Run a simple count query.
	q, _ := query.NewParser().ParseString("SELECT count() AS count")
	q.Prefix = req.FormValue("prefix")
	q.SetTable(table)
	q.SetFdb(s.fdb)

	return s.RunQuery(table, q)
}

// POST /tables/:name/query
func (s *Server) queryHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	q, err := s.parseQuery(table, params)
	if err != nil {
		return nil, err
	}

	return s.RunQuery(table, q)
}

// POST /tables/:name/query/codegen
func (s *Server) queryCodegenHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)

	// Return an error if the table already exists.
	table, err := s.OpenTable(vars["name"])
	if err != nil {
		return nil, err
	}

	q, err := s.parseQuery(table, params)
	if err != nil {
		return nil, err
	}

	// Generate the query source code.
	source, err := q.Codegen()
	return source, &TextPlainContentTypeError{}
}

func (s *Server) parseQuery(table *core.Table, params map[string]interface{}) (*query.Query, error) {
	var err error

	// DEPRECATED: Allow query to be passed in as root param.
	raw := params["query"]
	if raw == nil {
		raw = params
	}

	// Parse query if passed as a string. Otherwise deserialize JSON.
	var q *query.Query
	if str, ok := raw.(string); ok {
		if q, err = query.NewParser().ParseString(str); err != nil {
			return nil, err
		}
	} else if obj, ok := raw.(map[string]interface{}); ok {
		q = query.NewQuery()
		if err = q.Deserialize(obj); err != nil {
			return nil, err
		}
	}
	q.SetTable(table)
	q.SetFdb(s.fdb)

	return q, nil
}
