package skyd

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (s *Server) addEventHandlers(r *mux.Router) {
	r.HandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request) { s.getEventsHandler(w, req) }).Methods("GET")
	r.HandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request) { s.deleteEventsHandler(w, req) }).Methods("DELETE")

	r.HandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request) { s.getEventHandler(w, req) }).Methods("GET")
	r.HandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request) { s.replaceEventHandler(w, req) }).Methods("PUT")
	r.HandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request) { s.updateEventHandler(w, req) }).Methods("PATCH")
	r.HandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request) { s.deleteEventHandler(w, req) }).Methods("DELETE")
}

// GET /tables/:name/objects/:objectId/events
func (s *Server) getEventsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{}) (interface{}, error) {
		// Retrieve raw events.
		events, err := servlet.GetEvents(table, vars["objectId"])
		if err != nil {
			return nil, err
		}

		// Denormalize events.
		output := make([]map[string]interface{}, 0)
		for _, event := range events {
			e, err := table.SerializeEvent(event)
			if err != nil {
				return nil, err
			}
			output = append(output, e)
		}

		return output, nil
	})
}

// DELETE /tables/:name/objects/:objectId/events
func (s *Server) deleteEventsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{}) (interface{}, error) {
		return nil, servlet.DeleteEvents(table, vars["objectId"])
	})
}

// GET /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) getEventHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{}) (interface{}, error) {
		// Parse timestamp.
		timestamp, err := time.Parse(time.RFC3339, vars["timestamp"])
		if err != nil {
			return nil, err
		}

		// Find event.
		event, err := servlet.GetEvent(table, vars["objectId"], timestamp)
		if err != nil {
			return nil, err
		}
		if event == nil {
			return nil, errors.New("Event not found.")
		}

		// Convert an event to a serializable object.
		return table.SerializeEvent(event)
	})
}

// PUT /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) replaceEventHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{}) (interface{}, error) {
		params["timestamp"] = vars["timestamp"]
		event, err := table.DeserializeEvent(params)
		if err != nil {
			return nil, err
		}
		return nil, servlet.PutEvent(table, vars["objectId"], event, true)
	})
}

// PATCH /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) updateEventHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{}) (interface{}, error) {
		params["timestamp"] = vars["timestamp"]
		event, err := table.DeserializeEvent(params)
		if err != nil {
			return nil, err
		}
		return nil, servlet.PutEvent(table, vars["objectId"], event, false)
	})
}

// DELETE /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) deleteEventHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.processWithObject(w, req, vars["name"], vars["objectId"], func(table *Table, servlet *Servlet, params map[string]interface{}) (interface{}, error) {
		timestamp, err := time.Parse(time.RFC3339, vars["timestamp"])
		if err != nil {
			return nil, fmt.Errorf("Unable to parse timestamp: %v", timestamp)
		}

		return nil, servlet.DeleteEvent(table, vars["objectId"], timestamp)
	})
}
