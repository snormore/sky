package server

// eventHandler handles the management of events in the database.
type eventHandler struct {
	s *Server
}

// installEventHandler adds table routes to the server.
func installEventHandler(s *Server) *eventHandler {
	h := &eventHandler{s: s}

	s.HandleFunc("/tables/{table}/objects/{id}/events", EnsureTableHandler(HandleFunc(h.getEvents))).Methods("GET")
	s.HandleFunc("/tables/{table}/objects/{id}/events", EnsureTableHandler(HandleFunc(h.deleteEvents))).Methods("DELETE")

	s.HandleFunc("/tables/{table}/objects/{id}/events/{timestamp}", EnsureTableHandler(HandleFunc(h.getEvent))).Methods("GET")
	s.HandleFunc("/tables/{table}/objects/{id}/events/{timestamp}", EnsureMapHandler(EnsureTableHandler(HandleFunc(h.insertEvent)))).Methods("PUT", "PATCH")
	s.HandleFunc("/tables/{table}/objects/{id}/events/{timestamp}", EnsureTableHandler(HandleFunc(h.deleteEvent))).Methods("DELETE")

	s.HandleFunc("/events", HandleFunc(h.insertEventStream)).Methods("PUT", "PATCH")
	s.HandleFunc("/tables/{table}/events", EnsureTableHandler(HandleFunc(h.insertTableEventStream))).Methods("PUT", "PATCH")

	return h
}

// getEvents retrieves a list of all events associated with an object.
func (h *eventHandler) getEvents(s *Server, req Request) (interface{}, error) {
	warn("getEvents•1")
	t := req.Table()
	events, err := s.db.GetEvents(t.Name, req.Var("id"))
	if err != nil {
		return nil, err
	}
	if err := s.db.Factorizer().DefactorizeEvents(events, t.Name, t.PropertyFile()); err != nil {
		return nil, err
	}
	return t.SerializeEvents(events)
}

// deleteEvents deletes all events associated with an object.
func (h *eventHandler) deleteEvents(s *Server, req Request) (interface{}, error) {
	warn("deleteEvents•1")
	return nil, s.db.DeleteObject(req.Table().Name, req.Var("id"))
}

// getEvent retrieves a single event for an object at a given point in time.
func (h *eventHandler) getEvent(s *Server, req Request) (interface{}, error) {
	warn("getEvent•1")
	return nil, nil
}

// insertEvent adds a single event to an object.
func (h *eventHandler) insertEvent(s *Server, req Request) (interface{}, error) {
	t := req.Table()
	data := req.Data().(map[string]interface{})
	data["timestamp"] = req.Var("timestamp")

	event, err := t.DeserializeEvent(data)
	if err != nil {
		return nil, err
	}
	err = s.db.Factorizer().FactorizeEvent(event, t.Name, t.PropertyFile(), true)
	if err != nil {
		return nil, err
	}
	return nil, s.db.InsertEvent(t.Name, req.Var("id"), event)
}

// deleteEvent deletes a single event for an object at a given point in time.
func (h *eventHandler) deleteEvent(s *Server, req Request) (interface{}, error) {
	warn("deleteEvent•1")
	return nil, nil
}

// insertEventStream is a multi-table bulk insertion end point.
func (h *eventHandler) insertEventStream(s *Server, req Request) (interface{}, error) {
	return nil, nil
}

// insertTableEventStream is a single-table bulk insertion end point.
func (h *eventHandler) insertTableEventStream(s *Server, req Request) (interface{}, error) {
	return nil, nil
}
