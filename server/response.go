package server

// Response is a high level interface to the router methods on the server.
// It abstracts HTTP response to a more basic level.
type Response interface {
	Data() interface{}
	Error() error
}

// Response is the concrete implementation of the Response interface.
type response struct {
	data interface{}
	err  error
}

// Error returns the response error code.
func (r *response) Error() error {
	return r.err
}

// Data returns the parsed input data as a map.
func (r *response) Data() interface{} {
	return r.data
}
