package server

// queryHandler handles the execute of queries against database tables.
type queryHandler struct {
	s *Server
}

// installQueryHandler adds query routes to the server.
func installQueryHandler(s *Server) *queryHandler {
	h := &queryHandler{s: s}
	s.HandleFunc("/tables/{table}/query", EnsureTableHandler(HandleFunc(h.execute))).Methods("POST")
	return h
}

// execute reads the incoming query and executes it against the given table.
func (h *queryHandler) execute(s *Server, req Request) (interface{}, error) {
	// Open table.
	// t := req.Table()

	// TODO: Parse query.
	// TODO: Execute query.
	return nil, nil
}

