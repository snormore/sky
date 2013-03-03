package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addPropertyHandlers(r *mux.Router) {
	r.HandleFunc("/tables/{name}/properties", func(w http.ResponseWriter, req *http.Request) { s.getPropertiesHandler(w, req) }).Methods("GET")
	r.HandleFunc("/tables/{name}/properties", func(w http.ResponseWriter, req *http.Request) { s.createPropertyHandler(w, req) }).Methods("POST")
	r.HandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request) { s.getPropertyHandler(w, req) }).Methods("GET")
}

// GET /tables/:name/properties
func (s *Server) getPropertiesHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    // Return an error if the table already exists.
    vars := mux.Vars(req)
    table, err := s.OpenTable(vars["name"])
    if table == nil || err != nil {
      return nil, err
    }
    
    // Return all properties.
    return table.GetProperties()
  })
}

// POST /tables/:name/properties
func (s *Server) createPropertyHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    // Return an error if the table already exists.
    vars := mux.Vars(req)
    table, err := s.OpenTable(vars["name"])
    if table == nil || err != nil {
      return nil, err
    }
    
    // Add property to table.
    name, _ := params["name"].(string)
    typ, _ := params["type"].(string)
    dataType, _ := params["dataType"].(string)
    return table.CreateProperty(name, typ, dataType)
  })
}

// GET /tables/:name/properties/:propertyName
func (s *Server) getPropertyHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    // Return an error if the table already exists.
    vars := mux.Vars(req)
    table, err := s.OpenTable(vars["name"])
    if table == nil || err != nil {
      return nil, err
    }
    
    // Retrieve property.
    return table.GetPropertyByName(vars["propertyName"])
  })
}
