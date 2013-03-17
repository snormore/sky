package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addQueryHandlers() {
	s.ApiHandleFunc("/tables/{name}/query", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.queryHandler(w, req, params)
	}).Methods("POST")
}

// POST /tables/:name/query
func (s *Server) queryHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.RunQuery(vars["name"],
		"function aggregate(cursor, data)\n"+
		"  data.path_count = (data.path_count or 0) + 1\n"+
		"  while cursor:next() do\n"+
		"    data.event_count = (data.event_count or 0) + 1\n"+
		"  end\n"+
		"end\n" +
		"\n" +
		"function merge(results, data)\n" +
		"  results.path_count = (results.path_count or 0) + data.path_count\n"+
		"  results.event_count = (results.event_count or 0) + data.event_count\n"+
		"end\n",
	)
}
