package server

import (
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
	s.HandleFunc("/tables/{table}/query", EnsureTableHandler(EnsureMapHandler(HandleFunc(h.execute)))).Methods("POST")
	return h
}

// execute reads the incoming query and executes it against the given table.
func (h *queryHandler) execute(s *Server, req Request) (interface{}, error) {
	t := req.Table()
	data := req.Data().(map[string]interface{})
	querystring, _ := data["query"].(string)

	// Parse query.
	q, err := parser.ParseString(querystring)
	if err != nil {
		return nil, err
	}
	q.DynamicDecl = func (ref *ast.VarRef) *ast.VarDecl {
		p, _ := t.GetPropertyByName(ref.Name)
		if p == nil {
			return nil
		}
		return ast.NewVarDecl(p.Id, p.Name, p.DataType)
	}
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
		if err = mappers[i].Execute(cursor, "", result); err != nil {
			return nil, err
		}
		results[i] = result
	}
	mappers[0].Dump()

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
