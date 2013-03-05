package skyd

import (
  "errors"
  "fmt"
  "github.com/gorilla/mux"
  "net/http"
  "time"
)

func (s *Server) addEventHandlers(r *mux.Router) {
  r.HandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request) { s.replaceEventHandler(w, req) }).Methods("PUT")
}

// POST /tables/:name/properties
func (s *Server) replaceEventHandler(w http.ResponseWriter, req *http.Request) {
  vars := mux.Vars(req)
  s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{})(interface{}, error) {
    event := &Event{}
    
    // Parse timestamp.
    if timestamp, ok := params["timestamp"].(string); ok {
      ts, err := time.Parse(time.RFC3339, timestamp)
      if err != nil {
        return nil, fmt.Errorf("Unable to parse timestamp: %v", timestamp)
      }
      event.Timestamp = ts
    } else {
      return nil, errors.New("Timestamp required.")
    }
    
    // Convert maps to use property identifiers.
    if data, ok := params["data"].(map[interface{}]interface{}); ok {
      normalizedData, err := table.NormalizeMap(data)
      if err != nil {
        return nil, err
      }
      event.Data = normalizedData
    }
    if action, ok := params["action"].(map[interface{}]interface{}); ok {
      normalizedAction, err := table.NormalizeMap(action)
      if err != nil {
        return nil, err
      }
      event.Action = normalizedAction
    }
    
    return nil, servlet.PutEvent(table, vars["objectId"], event)
  })
}
