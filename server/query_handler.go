package server

import (
	"fmt"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/ast/validator"
	"github.com/skydb/sky/query/codegen/hashmap"
	"github.com/skydb/sky/query/codegen/mapper"
	"github.com/skydb/sky/query/codegen/reducer"
	"github.com/skydb/sky/query/parser"
)

// queryHandler handles the execute of queries against database tables.
type queryHandler struct {
	s *Server
}

// installQueryHandler adds query routes to the server.
func installQueryHandler(s *Server) *queryHandler {
	h := &queryHandler{s: s}
	s.HandleFunc("/tables/{table}/query", EnsureTableHandler(HandleFunc(h.query))).Methods("POST")
	s.HandleFunc("/tables/{table}/stats", EnsureTableHandler(HandleFunc(h.stats))).Methods("GET")
	return h
}

// query reads the incoming query and executes it against the given table.
func (h *queryHandler) query(s *Server, req Request) (interface{}, error) {
	var data map[string]interface{}
	switch d := req.Data().(type) {
	case map[string]interface{}:
		data = d
	case []byte:
		data = map[string]interface{}{"query": string(d)}
	default:
		return nil, fmt.Errorf("server: map or string input required")
	}
	querystring, _ := data["query"].(string)
	return h.execute(s, req, querystring)
}

// stats executes a simple event count against the given table.
func (h *queryHandler) stats(s *Server, req Request) (interface{}, error) {
	return h.execute(s, req, "SELECT count()")
}

// execute runs a query against the table.
func (h *queryHandler) execute(s *Server, req Request, querystring string) (interface{}, error) {
	t := req.Table()

	var data, ok = req.Data().(map[string]interface{})
	if !ok {
		data = make(map[string]interface{})
	}

	// Retrieve prefix.
	prefix, ok := data["prefix"].(string)
	if !ok {
		prefix = req.Var("prefix")
	}

	// Parse query.
	q, err := parser.ParseString(querystring)
	if err != nil {
		return nil, err
	}
	q.DynamicDecl = func(ref *ast.VarRef) *ast.VarDecl {
		p, _ := t.GetPropertyByName(ref.Name)
		if p == nil {
			return nil
		}
		return ast.NewVarDecl(p.Id, p.Name, p.DataType)
	}
	ast.Normalize(q)
	if err := q.Finalize(); err != nil {
		return nil, err
	}

	// Validate query.
	if err := validator.Validate(q); err != nil {
		return nil, err
	}

	// Retrieve factorizer and database cursors.
	f := s.db.TableFactorizer(t.Name)
	cursors, err := s.db.Cursors(t.Name)
	if err != nil {
		return nil, err
	}
	defer cursors.Close()

	// Execute one mapper for each cursor.
	mappers := make([]*mapper.Mapper, len(cursors))
	results := make([]*hashmap.Hashmap, len(cursors))
	for i := 0; i < len(mappers); i++ {
		cursor := cursors[i]

		var err error
		if mappers[i], err = mapper.New(q, f); err != nil {
			return nil, err
		}

		result := hashmap.New()
		if err = mappers[i].Execute(cursor, prefix, result); err != nil {
			return nil, err
		}
		results[i] = result
	}
	// mappers[0].Dump()

	// TODO: Run mappers in parallel.

	// Combine all the results into one final result.
	r := reducer.New(q, f)
	for _, result := range results {
		if err := r.Reduce(result); err != nil {
			return nil, err
		}
	}
	return r.Output(), nil
}
