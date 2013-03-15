package skyd

import (
  "fmt"
  "github.com/gorilla/mux"
  "github.com/jmhodges/levigo"
  "net/http"
)

func (s *Server) addQueryHandlers(r *mux.Router) {
  r.HandleFunc("/tables/{name}/query", func(w http.ResponseWriter, req *http.Request) { s.queryHandler(w, req) }).Methods("POST")
}

// POST /tables/:name/query
func (s *Server) queryHandler(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  s.processWithTable(w, req, vars["name"], func(table *Table, params map[string]interface{})(interface{}, error) {
    // TODO: Parse query from JSON POST.
    // TODO: Codegen Lua from query object.
    
    // TODO: Execute query on each servlet.
    // TODO: Merge results.
    // TODO: Return results to client.

    engine, err := NewExecutionEngine(table.propertyFile,
      "function aggregate(cursor, data)\n" +
      "  data.path_count = (data.path_count or 0) + 1\n" +
      "  while cursor:next() do\n" +
      "    data.event_count = (data.event_count or 0) + 1\n" +
      "  end\n" +
      "end\n",
    )
    if err != nil {
      return nil, err
    }

    // Initialize execution engine.
    err = engine.Init()
    if err != nil {
      return nil, err
    }
    defer engine.Destroy()

    // Initialize iterator.
    servlet := s.servlets[0]
    ro := levigo.NewReadOptions()
    defer ro.Close()
    iterator := servlet.db.NewIterator(ro)
    iterator.SeekToFirst()
    engine.SetIterator(iterator)

    // fmt.Println(engine.FullAnnotatedSource())

    // Run aggregation for the servlet.
    results, err := engine.Aggregate()

    return results, err
  })
}
