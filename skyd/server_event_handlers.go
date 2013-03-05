package skyd

import (
  "github.com/gorilla/mux"
  "net/http"
)

func (s *Server) addEventHandlers(r *mux.Router) {
  r.HandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request) { s.replaceEventHandler(w, req) }).Methods("PUT")
}

// POST /tables/:name/properties
func (s *Server) replaceEventHandler(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{})(interface{}, error) {
    event, err := table.DeserializeEvent(params)
    if err != nil {
      return nil, err
    }
    return nil, servlet.PutEvent(table, vars["objectId"], event)
  })
}
