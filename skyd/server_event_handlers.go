package skyd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

func (s *Server) addEventHandlers() {
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getEventsHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.deleteEventsHandler(w, req, params)
	}).Methods("DELETE")

	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.getEventHandler(w, req, params)
	}).Methods("GET")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.replaceEventHandler(w, req, params)
	}).Methods("PUT")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.updateEventHandler(w, req, params)
	}).Methods("PATCH")
	s.ApiHandleFunc("/tables/{name}/objects/{objectId}/events/{timestamp}", func(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
		return s.deleteEventHandler(w, req, params)
	}).Methods("DELETE")

	// Streaming import.
	s.router.HandleFunc("/tables/{name}/events", s.streamUpdateEventsHandler).Methods("PATCH")
}

// GET /tables/:name/objects/:objectId/events
func (s *Server) getEventsHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (interface{}, error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	// Retrieve raw events.
	events, _, err := servlet.GetEvents(table, vars["objectId"])
	if err != nil {
		return nil, err
	}

	// Denormalize events.
	output := make([]map[string]interface{}, 0)
	for _, event := range events {
		if err = table.DefactorizeEvent(event, s.factors); err != nil {
			return nil, err
		}
		e, err := table.SerializeEvent(event)
		if err != nil {
			return nil, err
		}
		output = append(output, e)
	}

	return output, nil
}

// DELETE /tables/:name/objects/:objectId/events
func (s *Server) deleteEventsHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	return nil, servlet.DeleteEvents(table, vars["objectId"])
}

// GET /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) getEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

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
	// Return an empty event if there isn't one.
	if event == nil {
		event = NewEvent(vars["timestamp"], map[int64]interface{}{})
	}

	// Convert an event to a serializable object.
	if err = table.DefactorizeEvent(event, s.factors); err != nil {
		return nil, err
	}
	e, err := table.SerializeEvent(event)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// PUT /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) replaceEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	params["timestamp"] = vars["timestamp"]
	event, err := table.DeserializeEvent(params)
	if err != nil {
		return nil, err
	}
	err = table.FactorizeEvent(event, s.factors, true)
	if err != nil {
		return nil, err
	}

	return nil, servlet.PutEvent(table, vars["objectId"], event, true)
}

// PATCH /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) updateEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	params["timestamp"] = vars["timestamp"]
	event, err := table.DeserializeEvent(params)
	if err != nil {
		return nil, err
	}
	err = table.FactorizeEvent(event, s.factors, true)
	if err != nil {
		return nil, err
	}
	return nil, servlet.PutEvent(table, vars["objectId"], event, false)
}

// PATCH /tables/:name/events
func (s *Server) streamUpdateEventsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	t0 := time.Now()

	events_written := 0
	err := func() error {
		// Stream in JSON event objects.
		decoder := json.NewDecoder(req.Body)
		for {
			// Read in a JSON object.
			rawEvent := map[string]interface{}{}
			if err := decoder.Decode(&rawEvent); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("Malformed json event: %v", err)
			}
			// Extract the object identifier.
			objectId, ok := rawEvent["id"].(string)
			if !ok {
				return fmt.Errorf("Object identifier required")
			}

			// Determine the appropriate servlet to insert into.
			table, servlet, err := s.GetObjectContext(vars["name"], objectId)
			if err != nil {
				return fmt.Errorf("Cannot determine object context: %v", err)
			}

			// Convert to a Sky event and insert.
			event, err := table.DeserializeEvent(rawEvent)
			if err != nil {
				return fmt.Errorf("Cannot deserialize: %v", err)
			}
			if err = table.FactorizeEvent(event, s.factors, true); err != nil {
				return fmt.Errorf("Cannot factorize: %v", err)
			}
			if err = servlet.PutEvent(table, objectId, event, false); err != nil {
				return fmt.Errorf("Cannot put event: %v", err)
			}
			events_written++
		}
		return nil
	}()

	if err != nil {
		s.logger.Printf("ERR %v", err)
		fmt.Fprintf(w, `{"message":"%v"}`, err)
	}
	s.logger.Printf("%s \"%s %s %s %d events OK\" %0.3f", req.RemoteAddr, req.Method, req.URL.Path, req.Proto, events_written, time.Since(t0).Seconds())
}

// DELETE /tables/:name/objects/:objectId/events/:timestamp
func (s *Server) deleteEventHandler(w http.ResponseWriter, req *http.Request, params map[string]interface{}) (ret interface{}, err error) {
	vars := mux.Vars(req)
	table, servlet, err := s.GetObjectContext(vars["name"], vars["objectId"])
	if err != nil {
		return nil, err
	}

	timestamp, err := time.Parse(time.RFC3339, vars["timestamp"])
	if err != nil {
		return nil, fmt.Errorf("Unable to parse timestamp: %v", timestamp)
	}

	return nil, servlet.DeleteEvent(table, vars["objectId"], timestamp)
}
