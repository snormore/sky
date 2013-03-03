package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addPropertyHandlers(r *mux.Router) {
	r.HandleFunc("/tables/{name}/properties", func(w http.ResponseWriter, req *http.Request) { s.getPropertiesHandler(w, req) }).Methods("GET")
	r.HandleFunc("/tables/{name}/properties", func(w http.ResponseWriter, req *http.Request) { s.createPropertyHandler(w, req) }).Methods("POST")
	r.HandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request) { s.getPropertyHandler(w, req) }).Methods("GET")
	r.HandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request) { s.updatePropertyHandler(w, req) }).Methods("PATCH")
	r.HandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request) { s.deletePropertyHandler(w, req) }).Methods("DELETE")
}

// GET /tables/:name/properties
func (s *Server) getPropertiesHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
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
    vars := mux.Vars(req)
    table, err := s.OpenTable(vars["name"])
    if table == nil || err != nil {
      return nil, err
    }
    
    // Retrieve property.
    return table.GetPropertyByName(vars["propertyName"])
  })
}

// PATCH /tables/:name/properties/:propertyName
func (s *Server) updatePropertyHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    vars := mux.Vars(req)
    table, err := s.OpenTable(vars["name"])
    if table == nil || err != nil {
      return nil, err
    }
    
    // Retrieve property.
    property, err := table.GetPropertyByName(vars["propertyName"])
    if err != nil {
      return nil, err
    }

    // Update property and save property file.
    name, _ := params["name"].(string)
    property.Name = name
    err = table.SavePropertyFile()
    if err != nil {
      return nil, err
    }
    
    return property, nil
  })
}

// DELETE /tables/:name/properties/:propertyName
func (s *Server) deletePropertyHandler(w http.ResponseWriter, req *http.Request) {
  s.process(w, req, func(params map[string]interface{})(interface{}, error) {
    vars := mux.Vars(req)
    table, err := s.OpenTable(vars["name"])
    if table == nil || err != nil {
      return nil, err
    }
    
    // Retrieve property.
    property, err := table.GetPropertyByName(vars["propertyName"])
    if err != nil {
      return nil, err
    }

    // Delete property and save property file.
    table.DeleteProperty(property)
    err = table.SavePropertyFile()
    if err != nil {
      return nil, err
    }
    
    return nil, nil
  })
}
