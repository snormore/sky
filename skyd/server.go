package skyd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmhodges/levigo"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	DefaultPort = 8585

	ServerMessageTypeExecute  = "execute"
	ServerMessageTypeShutdown = "shutdown"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Server is the front end that controls access to tables.
type Server struct {
	httpServer *http.Server
	router     *mux.Router
	path       string
	listener   net.Listener
	servlets   []*Servlet
	tables     map[string]*Table
	channel    chan *Message
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewServer returns a new Server.
func NewServer(port uint, path string) *Server {
	r := mux.NewRouter()
	s := &Server{
		httpServer: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r},
		router:     r,
		path:       path,
		tables:     make(map[string]*Table),
		channel:    make(chan *Message),
	}

	s.addHandlers()
	s.addTableHandlers()
	s.addPropertyHandlers()
	s.addEventHandlers()
	s.addQueryHandlers()

	return s
}

//------------------------------------------------------------------------------
//
// Properties
//
//------------------------------------------------------------------------------

// The root server path.
func (s *Server) Path() string {
	return s.path
}

// The path to the data directory.
func (s *Server) DataPath() string {
	return fmt.Sprintf("%v/data", s.path)
}

// The path to the table metadata directory.
func (s *Server) TablesPath() string {
	return fmt.Sprintf("%v/tables", s.path)
}

// Generates the path for a table attached to the server.
func (s *Server) TablePath(name string) string {
	return fmt.Sprintf("%v/%v", s.TablesPath(), name)
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Lifecycle
//--------------------------------------

// Runs the server.
func (s *Server) ListenAndServe() error {
	err := s.open()
	if err != nil {
		fmt.Printf("Unable to open server: %v", err)
		return err
	}

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.close()
		return err
	}
	s.listener = listener
	go s.process()
	go s.httpServer.Serve(s.listener)
	return nil
}

// Stops the server.
func (s *Server) Shutdown() error {
	// Close servlets.
	s.close()

	if s.listener != nil {
		// Wait for the message loop to shutdown.
		m := NewShutdownMessage()
		s.channel <- m
		m.wait()

		// Cleanup channel.
		close(s.channel)
		s.channel = nil

		// Then stop the server.
		err := s.listener.Close()
		s.listener = nil
		return err
	}
	return nil
}

// Checks if the server is listening for new connections.
func (s *Server) Running() bool {
	return (s.listener != nil)
}

// Opens the data directory and servlets.
func (s *Server) open() error {
	s.close()

	// Setup the file system if it doesn't exist.
	err := s.createIfNotExists()
	if err != nil {
		panic(fmt.Sprintf("skyd.Server: Unable to create server folders: %v", err))
	}

	// Create servlets from child directories with numeric names.
	infos, err := ioutil.ReadDir(s.DataPath())
	if err != nil {
		return err
	}
	for _, info := range infos {
		match, _ := regexp.MatchString("^\\d$", info.Name())
		if info.IsDir() && match {
			s.servlets = append(s.servlets, NewServlet(fmt.Sprintf("%s/%s", s.DataPath(), info.Name())))
		}
	}

	// If none exist then build them based on the number of logical CPUs available.
	if len(s.servlets) == 0 {
		cpuCount := runtime.NumCPU()
		cpuCount = 1
		for i := 0; i < cpuCount; i++ {
			s.servlets = append(s.servlets, NewServlet(fmt.Sprintf("%s/%v", s.DataPath(), i)))
		}
	}

	// Open servlets.
	for _, servlet := range s.servlets {
		err = servlet.Open()
		if err != nil {
			s.close()
			return err
		}
	}

	return nil
}

// Closes the data directory and servlets.
func (s *Server) close() {
	// Close servlets.
	if s.servlets != nil {
		for _, servlet := range s.servlets {
			servlet.Close()
		}
		s.servlets = nil
	}
}

// Creates the appropriate directory structure if one does not exist.
func (s *Server) createIfNotExists() error {
	// Create root directory.
	err := os.MkdirAll(s.path, 0700)
	if err != nil {
		return err
	}

	// Create data directory and one directory for each servlet.
	err = os.MkdirAll(s.DataPath(), 0700)
	if err != nil {
		return err
	}

	// Create tables directory.
	err = os.MkdirAll(s.TablesPath(), 0700)
	if err != nil {
		return err
	}

	return nil
}

//--------------------------------------
// Routing
//--------------------------------------

// Serially processes server messages routed through the server channel.
func (s *Server) ApiHandleFunc(route string, handlerFunction func(http.ResponseWriter, *http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunction := func(w http.ResponseWriter, req *http.Request) {
		var ret interface{}
		params, err := s.decodeParams(w, req)
		if err == nil {
			ret, err = handlerFunction(w, req, params)
		}

		// If there is an error then replace the return value.
		if err != nil {
			ret = map[string]interface{}{"message": err.Error()}
		}

		// Write header status.
		if err == nil {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Encode the return value appropriately.
		if ret != nil {
			contentType := req.Header.Get("Content-Type")
			if contentType == "application/json" {
				encoder := json.NewEncoder(w)
				err := encoder.Encode(ConvertToStringKeys(ret))
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}

	return s.router.HandleFunc(route, wrappedFunction)
}

//--------------------------------------
// Synchronization
//--------------------------------------

// Executes a function through a single-threaded server context.
func (s *Server) sync(f func() (interface{}, error)) (interface{}, error) {
	m := s.async(f)
	return m.wait()
}

// Executes a function through a single-threaded server context.
func (s *Server) async(f func() (interface{}, error)) *Message {
	m := NewExecuteMessage(f)
	s.channel <- m
	return m
}

// Executes a function through a single-threaded server context.
func (s *Server) executeWithTable(tableName string, f func(table *Table) (interface{}, error)) (interface{}, error) {
	return s.sync(func() (interface{}, error) {
		// Return an error if the table already exists.
		table, err := s.OpenTable(tableName)
		if err != nil {
			return nil, err
		}
		return f(table)
	})
}

// Executes a function through a single-threaded servlet context.
func (s *Server) executeWithObject(tableName string, objectId string, f func(servlet *Servlet, table *Table) (interface{}, error)) (interface{}, error) {
	var m *Message
	_, err := s.sync(func() (interface{}, error) {
		// Return an error if the table already exists.
		table, err := s.OpenTable(tableName)
		if err != nil {
			return nil, err
		}

		// Determine servlet index.
		index, err := s.GetObjectServletIndex(table, objectId)
		if err != nil {
			return nil, err
		}
		servlet := s.servlets[index]

		// Pass off the table to the servlet message loop.
		m = servlet.async(func() (interface{}, error) {
			return f(servlet, table)
		})

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	// Make sure we've exited the server loop and we're waiting on the servlet loop.
	return m.wait()
}

//--------------------------------------
// Message Loop
//--------------------------------------

// Serially processes server messages routed through the server channel.
func (s *Server) process() {
	for message := range s.channel {
		message.execute()
		if message.messageType == ShutdownMessageType {
			return
		}
	}
}

// Decodes the body of the message into parameters.
func (s *Server) decodeParams(w http.ResponseWriter, req *http.Request) (map[string]interface{}, error) {
	// Parses body parameters.
	params := make(map[string]interface{})
	contentType := req.Header.Get("Content-Type")
	if contentType == "application/json" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&params)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusBadRequest)
			return nil, errors.New("Invalid body.")
		}
	} else {
		w.WriteHeader(http.StatusNotAcceptable)
		return nil, errors.New("Invalid content type.")
	}
	return params, nil
}

//--------------------------------------
// Servlet Management
//--------------------------------------

// Calculates a tablet index based on the object identifier even hash.
func (s *Server) GetObjectServletIndex(t *Table, objectId string) (uint32, error) {
	// Encode object identifier.
	encodedObjectId, err := t.EncodeObjectId(objectId)
	if err != nil {
		return 0, err
	}

	// Calculate the even bits of the FNV1a hash.
	h := fnv.New64a()
	h.Reset()
	h.Write(encodedObjectId)
	hashcode := h.Sum64()
	index := CondenseUint64Even(hashcode) % uint32(len(s.servlets))

	return index, nil
}

//--------------------------------------
// Table Management
//--------------------------------------

// Retrieves a table that has already been opened.
func (s *Server) GetTable(name string) *Table {
	return s.tables[name]
}

// Opens a table and returns a reference to it.
func (s *Server) OpenTable(name string) (*Table, error) {
	// If table already exists then use it.
	table := s.GetTable(name)
	if table != nil {
		return table, nil
	}

	// Otherwise open it and save the reference.
	table = NewTable(name, s.TablePath(name))
	err := table.Open()
	if err != nil {
		table.Close()
		return nil, err
	}
	s.tables[name] = table

	return table, nil
}

//--------------------------------------
// Query
//--------------------------------------

// Runs a query against a table.
func (s *Server) RunQuery(tableName string, json map[string]interface{}) (interface{}, error) {
	var engine *ExecutionEngine

	// Create a channel to receive aggregate responses.
	rchannel := make(chan *Message, len(s.servlets))

	// Retrieve table and setup servlets within the server context.
	_, err := s.sync(func() (interface{}, error) {
		// Return an error if the table already exists.
		table, err := s.OpenTable(tableName)
		if err != nil {
			return nil, err
		}

		// Deserialize the query.
		query := NewQuery(table)
		err = query.Deserialize(json)
		if err != nil {
			return nil, err
		}
		
		// Generate the query source code.
		source, err := query.Codegen()
		if err != nil {
			return nil, err
		}
		
		// Create an engine for merging results.
		engine, err = NewExecutionEngine(table.propertyFile, source)
		if err != nil {
			return nil, err
		}
		// fmt.Println(engine.FullAnnotatedSource())

		// Execute across each server context.
		for _, servlet := range s.servlets {

			// Execute servlets asynchronously and retrieve responses outside
			// of the server context.
			rchannel <- servlet.async(func() (interface{}, error) {
				// Create an engine for each servlet.
				e, err := NewExecutionEngine(table.propertyFile, source)
				if err != nil {
					return nil, err
				}
				defer e.Destroy()

				// Initialize iterator.
				ro := levigo.NewReadOptions()
				defer ro.Close()
				iterator := servlet.db.NewIterator(ro)
				iterator.SeekToFirst()
				e.SetIterator(iterator)

				// Run aggregation for the servlet.
				return e.Aggregate()
			})
		}

		return nil, nil
	})
	if engine != nil {
		defer engine.Destroy()
	}
	if err != nil {
		return nil, err
	}

	// Wait for each servlet to complete and then merge the results.
	var servletError error
	var result interface{}
	result = make(map[interface{}]interface{})
	for {
		var m *Message
		select {
		case m = <-rchannel:
		default:
			m = nil
		}
		if m == nil {
			break
		}
		ret, err := m.wait()
		if err != nil {
			fmt.Printf("skyd.Server: Aggregate error: %v", err)
			servletError = err
		}
		if ret != nil {
			result, err = engine.Merge(result, ret)
			if err != nil {
				fmt.Printf("skyd.Server: Merge error: %v", err)
				servletError = err
			}
		}
	}
	err = servletError
	return result, err
}
