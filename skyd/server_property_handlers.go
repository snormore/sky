package skyd

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) addPropertyHandlers() {
	s.ApiHandleFunc("/tables/{name}/properties", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getPropertiesHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/properties", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.createPropertyHandler(w, req, params)
	}).Methods("POST")

	s.ApiHandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getPropertyHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.updatePropertyHandler(w, req, params)
	}).Methods("PATCH")
	s.ApiHandleFunc("/tables/{name}/properties/{propertyName}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.deletePropertyHandler(w, req, params)
	}).Methods("DELETE")
}

// GET /tables/:name/properties
func (s *Server) getPropertiesHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.executeWithTable(vars["name"], func(table *Table) (interface{}, error) {
		return table.GetProperties()
	})
}

// POST /tables/:name/properties
func (s *Server) createPropertyHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.executeWithTable(vars["name"], func(table *Table) (interface{}, error) {
		name, _ := params["name"].(string)
		typ, _ := params["type"].(string)
		dataType, _ := params["dataType"].(string)
		return table.CreateProperty(name, typ, dataType)
	})
}

// GET /tables/:name/properties/:propertyName
func (s *Server) getPropertyHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.executeWithTable(vars["name"], func(table *Table) (interface{}, error) {
		return table.GetPropertyByName(vars["propertyName"])
	})
}

// PATCH /tables/:name/properties/:propertyName
func (s *Server) updatePropertyHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.executeWithTable(vars["name"], func(table *Table) (interface{}, error) {
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
func (s *Server) deletePropertyHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	return s.executeWithTable(vars["name"], func(table *Table) (interface{}, error) {
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
