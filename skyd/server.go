package skyd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/benbjohnson/go-raft"
	"github.com/gorilla/mux"
	"hash/fnv"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"sync"
	"time"
)

//------------------------------------------------------------------------------
//
// Globals
//
//------------------------------------------------------------------------------

const (
	DefaultElectionTimeout  = 2 * time.Second
	DefaultHeartbeatTimeout = 1 * time.Second
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Server is the front end that controls access to tables.
type Server struct {
	httpServer       *http.Server
	raftServer       *raft.Server
	cluster          *Cluster
	router           *mux.Router
	logger           *log.Logger
	path             string
	host             string
	port             uint
	listener         net.Listener
	servlets         []*Servlet
	tables           map[string]*Table
	factors          *Factors
	shutdownChannel  chan bool
	shutdownFinished chan bool
	mutex            sync.Mutex
	ElectionTimeout  time.Duration
	HeartbeatTimeout time.Duration
}

//------------------------------------------------------------------------------
//
// Errors
//
//------------------------------------------------------------------------------

//--------------------------------------
// Alternate Content Type
//--------------------------------------

// Hackish way to return plain text for query codegen. Might fix later but
// it's rarely used.
type TextPlainContentTypeError struct {
}

func (e *TextPlainContentTypeError) Error() string {
	return ""
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewServer returns a new Server.
func NewServer(port uint, path string) (*Server, error) {
	// Retrieve the host information.
	host, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("skyd.Server: Unable to determine hostname: %s", err)
	}

	r := mux.NewRouter()
	s := &Server{
		httpServer:       &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r},
		router:           r,
		logger:           log.New(os.Stdout, "", log.LstdFlags),
		path:             path,
		host:             host,
		port:             port,
		tables:           make(map[string]*Table),
		ElectionTimeout:  DefaultElectionTimeout,
		HeartbeatTimeout: DefaultHeartbeatTimeout,
	}

	s.addHandlers()
	s.addTableHandlers()
	s.addPropertyHandlers()
	s.addEventHandlers()
	s.addQueryHandlers()
	s.addClusterHandlers()
	s.addDebugHandlers()

	return s, nil
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

// The path to the factors database.
func (s *Server) FactorsPath() string {
	return fmt.Sprintf("%v/factors", s.path)
}

// The path to the binlog directory. This is used to log updates to the server
// for the purpose of replication.
func (s *Server) BinlogPath() string {
	return fmt.Sprintf("%v/binlog", s.path)
}

// The path to the raft directory. This is used for consensus information
// for replicating server configuration.
func (s *Server) RaftPath() string {
	return fmt.Sprintf("%v/raft", s.path)
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
func (s *Server) ListenAndServe(shutdownChannel chan bool) error {
	s.shutdownChannel = shutdownChannel

	err := s.open()
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		s.close()
		return err
	}
	s.listener = listener

	s.shutdownFinished = make(chan bool)
	go func() {
		s.httpServer.Serve(s.listener)
		s.shutdownFinished <- true
	}()

	s.logger.Printf("Sky v%s is now listening on http://localhost%s\n", Version, s.httpServer.Addr)

	return nil
}

// Stops the server.
func (s *Server) Shutdown() error {
	// Close servlets.
	s.close()

	// Close socket.
	if s.listener != nil {
		// Then stop the server.
		err := s.listener.Close()
		s.listener = nil
		if err != nil {
			return err
		}
	}
	// wait for server goroutine to finish
	<-s.shutdownFinished

	// Notify that the server is shutdown.
	if s.shutdownChannel != nil {
		s.shutdownChannel <- true
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
		return fmt.Errorf("skyd.Server: Unable to create server folders: %v", err)
	}

	// Open factors database.
	s.factors = NewFactors(s.FactorsPath())
	err = s.factors.Open()
	if err != nil {
		s.close()
		return err
	}

	// Create servlets from child directories with numeric names.
	infos, err := ioutil.ReadDir(s.DataPath())
	if err != nil {
		return err
	}
	for _, info := range infos {
		match, _ := regexp.MatchString("^\\d$", info.Name())
		if info.IsDir() && match {
			s.servlets = append(s.servlets, NewServlet(fmt.Sprintf("%s/%s", s.DataPath(), info.Name()), s.factors))
		}
	}

	// If none exist then build them based on the number of logical CPUs available.
	if len(s.servlets) == 0 {
		cpuCount := runtime.NumCPU()
		for i := 0; i < cpuCount; i++ {
			s.servlets = append(s.servlets, NewServlet(fmt.Sprintf("%s/%v", s.DataPath(), i), s.factors))
		}
	}

	// Initialize empty cluster.
	s.cluster = NewCluster()

	// Initialize consensus server.
	if err = s.initRaftServer(); err != nil {
		s.close()
		return err
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

// Sets up and starts the consensus server.
func (s *Server) initRaftServer() error {
	// Setup the consensus server.
	var err error
	if s.raftServer, err = raft.NewServer(NewNodeId(), s.RaftPath(), s, s); err != nil {
		return err
	}
	s.raftServer.SetElectionTimeout(s.ElectionTimeout)
	s.raftServer.SetHeartbeatTimeout(s.HeartbeatTimeout)

	// Start the consensus server.
	if err = s.raftServer.Start(); err != nil {
		return err
	}

	// If this is a new server then add a group and a single node.
	if s.raftServer.IsLogEmpty() {
		s.raftServer.Initialize()
		nodeGroupId := NewNodeGroupId()
		if err = s.raftServer.Do(&CreateNodeGroupCommand{NodeGroupId: nodeGroupId}); err != nil {
			return err
		}
		if err = s.raftServer.Do(NewCreateNodeCommand(NewNodeId(), nodeGroupId, s.host, s.port)); err != nil {
			return err
		}
	}

	return nil
}

// Closes the data directory and servlets.
func (s *Server) close() {
	// Close consensus server.
	if s.raftServer != nil {
		s.raftServer.Stop()
		s.raftServer = nil
	}

	// Remove cluster info.
	s.cluster = nil

	// Close servlets.
	if s.servlets != nil {
		for _, servlet := range s.servlets {
			servlet.Close()
		}
		s.servlets = nil
	}

	// Close factors database.
	if s.factors != nil {
		s.factors.Close()
		s.factors = nil
	}
}

// Creates the appropriate directory structure if one does not exist.
func (s *Server) createIfNotExists() error {
	if err := os.MkdirAll(s.path, 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(s.DataPath(), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(s.TablesPath(), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(s.RaftPath(), 0700); err != nil {
		return err
	}
	return nil
}

//--------------------------------------
// Logging
//--------------------------------------

// Silences the log.
func (s *Server) Silence() {
	s.logger = log.New(ioutil.Discard, "", log.LstdFlags)
}

//--------------------------------------
// Routing
//--------------------------------------

// Parses incoming JSON objects and converts outgoing responses to JSON.
func (s *Server) ApiHandleFunc(route string, handlerFunction func(http.ResponseWriter, *http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	return s.apiHandleFunc(route, handlerFunction)
}

func (s *Server) apiHandleFunc(route string, handlerFunction func(http.ResponseWriter, *http.Request, map[string]interface{}) (interface{}, error)) *mux.Route {
	wrappedFunction := func(w http.ResponseWriter, req *http.Request) {
		// warn("%s \"%s %s %s\"", req.RemoteAddr, req.Method, req.RequestURI, req.Proto)
		t0 := time.Now()

		var ret interface{}
		params, err := s.decodeParams(w, req)
		if err == nil {
			ret, err = handlerFunction(w, req, params)
		}

		// If we're returning plain text then just dump out what's returned.
		if _, ok := err.(*TextPlainContentTypeError); ok {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			if str, ok := ret.(string); ok {
				w.Write([]byte(str))
			}
			return
		}

		// If there is an error then replace the return value.
		if err != nil {
			ret = map[string]interface{}{"message": err.Error()}
		}

		// Write header status.
		w.Header().Set("Content-Type", "application/json")
		var status int
		if err == nil {
			status = http.StatusOK
		} else {
			status = http.StatusInternalServerError
		}
		w.WriteHeader(status)

		// Write to access log.
		s.logger.Printf("%s \"%s %s %s\" %d %0.3f", req.RemoteAddr, req.Method, req.RequestURI, req.Proto, status, time.Since(t0).Seconds())
		if status != http.StatusOK {
			s.logger.Printf("ERROR %v", err)
		}

		// Encode the return value appropriately.
		if ret != nil {
			encoder := json.NewEncoder(w)
			err := encoder.Encode(ConvertToStringKeys(ret))
			if err != nil {
				fmt.Printf("skyd.Server: Encoding error: %v\n", err.Error())
				return
			}
		}
	}

	return s.router.HandleFunc(route, wrappedFunction)
}

// Decodes the body of the message into parameters.
func (s *Server) decodeParams(w http.ResponseWriter, req *http.Request) (map[string]interface{}, error) {
	// Parses body parameters.
	params := make(map[string]interface{})
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("Malformed json request: %v", err)
	}
	return params, nil
}

//--------------------------------------
// Servlet Management
//--------------------------------------

// Executes a function through a single-threaded servlet context.
func (s *Server) GetObjectContext(tableName string, objectId string) (*Table, *Servlet, error) {
	// Return an error if the table already exists.
	table, err := s.OpenTable(tableName)
	if err != nil {
		return nil, nil, err
	}

	// Determine servlet index.
	index, err := s.GetObjectServletIndex(table, objectId)
	if err != nil {
		return nil, nil, err
	}
	servlet := s.servlets[index]

	return table, servlet, nil
}

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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.tables[name]
}

// Retrieves a list of all tables in the database but does not open them.
// Do not use these table references for anything but informational purposes!
func (s *Server) GetAllTables() ([]*Table, error) {
	// Create a table object for each directory in the tables path.
	infos, err := ioutil.ReadDir(s.TablesPath())
	if err != nil {
		return nil, err
	}

	tables := []*Table{}
	for _, info := range infos {
		if info.IsDir() {
			tables = append(tables, NewTable(info.Name(), s.TablePath(info.Name())))
		}
	}

	return tables, nil
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

	s.mutex.Lock()
	s.tables[name] = table
	s.mutex.Unlock()

	return table, nil
}

// Deletes a table.
func (s *Server) DeleteTable(name string) error {
	// Return an error if the table doesn't exist.
	table := s.GetTable(name)
	if table == nil {
		table = NewTable(name, s.TablePath(name))
	}
	if !table.Exists() {
		return fmt.Errorf("Table does not exist: %s", name)
	}

	// Delete data from each servlet.
	for _, servlet := range s.servlets {
		if err := servlet.DeleteTable(name); err != nil {
			return err
		}
	}

	// Remove the table from the lookup and remove it's schema.
	s.mutex.Lock()
	delete(s.tables, name)
	defer s.mutex.Unlock()

	return table.Delete()
}

//--------------------------------------
// Query
//--------------------------------------

// Runs a query against a table.
func (s *Server) RunQuery(table *Table, query *Query) (interface{}, error) {
	var engine *ExecutionEngine
	engines := make([]*ExecutionEngine, 0)

	// Create a channel to receive aggregate responses.
	rchannel := make(chan interface{}, len(s.servlets))

	// Generate the query source code.
	source, err := query.Codegen()
	if err != nil {
		return nil, err
	}

	// Create an engine for merging results.
	engine, err = NewExecutionEngine(table, query.Prefix, source)
	if err != nil {
		return nil, err
	}
	defer engine.Destroy()
	//fmt.Println(engine.FullAnnotatedSource())

	// Initialize one execution engine for each servlet.
	var data interface{}
	for index, servlet := range s.servlets {
		// Create an engine for each servlet. The execution engine is
		// protected by a mutex so it's safe to destroy it at any time.
		e, err := servlet.CreateExecutionEngine(table, query.Prefix, source)
		if err != nil {
			return nil, err
		}
		defer e.Destroy()
		engines = append(engines, e)

		// Run initialization once if required.
		if index == 0 && query.RequiresInitialization() {
			if data, err = e.Initialize(); err != nil {
				return nil, err
			}
			// Reset the iterator.
			if err = e.ResetLmdbCursor(); err != nil {
				return nil, err
			}
		}
	}

	// Execute servlets asynchronously and retrieve responses outside
	// of the server context.
	for index, _ := range s.servlets {
		e := engines[index]
		go func() {
			if result, err := e.Aggregate(data); err != nil {
				rchannel <- err
			} else {
				rchannel <- result
			}
		}()
	}

	// Wait for each servlet to complete and then merge the results.
	var servletError error
	var result interface{}
	result = make(map[interface{}]interface{})
	for i := 0; i < len(s.servlets); i++ {
		ret := <-rchannel
		if err, ok := ret.(error); ok {
			fmt.Printf("skyd.Server: Aggregate error: %v", err)
			servletError = err
		} else {
			// Defactorize aggregate results.
			if err = query.Defactorize(ret); err != nil {
				return nil, err
			}

			// Merge results.
			if ret != nil {
				if result, err = engine.Merge(result, ret); err != nil {
					fmt.Printf("skyd.Server: Merge error: %v", err)
					servletError = err
				}
			}
		}
	}
	err = servletError

	return result, err
}

//--------------------------------------
// Raft Execution
//--------------------------------------

// Executes a command on the server if the server is the leader. Otherwise
// it forwards the command to the leader.
func (s *Server) Do(command raft.Command) error {
	warn("do.1 %v", command)
	if s.raftServer.State() == raft.Follower {
		
	}
	// TODO: Check if this is the leader.
	// TODO: If not then forward the command.
	return nil
}

//--------------------------------------
// Raft Transport
//--------------------------------------

// Sends AppendEntries RPCs to a peer when the server is the leader.
func (s *Server) SendAppendEntriesRequest(raftServer *raft.Server, peer *raft.Peer, req *raft.AppendEntriesRequest) (*raft.AppendEntriesResponse, error) {
	// Find connection info.
	host, port, err := s.cluster.GetNodeHostname(peer.Name())
	if err != nil {
		return nil, fmt.Errorf("skyd.Server: Peer not found: %s", peer.Name())
	}

	// Send the entries over HTTP.
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(req)
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/raft/append", host, port), "application/json", &b)
	if resp == nil {
		return nil, fmt.Errorf("skyd.Server: Append entries connection failed: %v", err)
	}
	defer resp.Body.Close()

	// Process the response.
	r := &raft.AppendEntriesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil && err != io.EOF {
		return nil, err
	}
	return r, nil
}

// Sends RequestVote RPCs to a peer when the server is the candidate.
func (s *Server) SendVoteRequest(raftServer *raft.Server, peer *raft.Peer, req *raft.RequestVoteRequest) (*raft.RequestVoteResponse, error) {
	// Find connection info.
	host, port, err := s.cluster.GetNodeHostname(peer.Name())
	if err != nil {
		return nil, fmt.Errorf("skyd.Server: Peer not found: %s", peer.Name())
	}

	// Send the entries over HTTP.
	var b bytes.Buffer
	json.NewEncoder(&b).Encode(req)
	resp, err := http.Post(fmt.Sprintf("http://%s:%d/raft/vote", host, port), "application/json", &b)
	if resp == nil {
		return nil, fmt.Errorf("skyd.Server: Request vote connection failed: %v", err)
	}
	defer resp.Body.Close()

	// Process the response.
	r := &raft.RequestVoteResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil && err != io.EOF {
		return nil, err
	}
	return r, nil
}
