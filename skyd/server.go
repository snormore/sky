package skyd

import (
  "encoding/json"
  "errors"
  "fmt"
  "github.com/gorilla/mux"
  "hash/fnv"
  "io"
  "io/ioutil"
  "net"
  "net/http"
  "os"
  "regexp"
  "runtime"
)

const (
  DefaultPort = 8585
)

// A Server is the front end that controls access to tables.
type Server struct {
  httpServer *http.Server
  path       string
  listener   net.Listener
  servlets   []*Servlet
  tables     map[string]*Table
  channel    chan *ServerMessage
}

// The request, response and handler function used for processing a message.
type ServerMessage struct {
  w       http.ResponseWriter
  req     *http.Request
  params  map[string]interface{}
  f       func(params map[string]interface{})(interface{}, error)
  channel chan bool
}

// NewServer returns a new Server.
func NewServer(port uint, path string) *Server {
  r := mux.NewRouter()
  s := &Server{
    httpServer: &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r},
    path:       path,
    tables:     make(map[string]*Table),
    channel:    make(chan *ServerMessage),
  }

  s.addHandlers(r)
  s.addTableHandlers(r)
  s.addPropertyHandlers(r)

  return s
}

func NewServerMessage(w http.ResponseWriter, req *http.Request, params map[string]interface{}, f func(params map[string]interface{})(interface{}, error)) *ServerMessage {
  return &ServerMessage{w:w, req:req, params:params, f:f, channel:make(chan bool)}
}

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
  go s.processMessages()
  go s.httpServer.Serve(s.listener)
  return nil
}

// Stops the server.
func (s *Server) Shutdown() error {
  if s.listener != nil {
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
    for i := 0; i < runtime.NumCPU(); i++ {
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

// Processes a request and return the appropriate data format.
func (s *Server) process(w http.ResponseWriter, req *http.Request, f func(map[string]interface{})(interface{}, error)) {
  params, err := decodeParams(w, req)
  if err != nil {
    return
  }

  // Push the message onto a queue to be processed serially.
  m := NewServerMessage(w, req, params, f)
  s.channel <- m
  <- m.channel
}

// Processes a request and automatically opens a given table.
func (s *Server) processWithTable(w http.ResponseWriter, req *http.Request, tableName string, f func(*Table, map[string]interface{})(interface{}, error)) {
  params, err := decodeParams(w, req)
  if err != nil {
    return
  }

  // Create a wrapper function to open the table.
  preprocess := func(_ map[string]interface{})(interface{}, error) {
    table, err := s.OpenTable(tableName)
    if table == nil || err != nil {
      return nil, err
    }
    
    // Execute the original function.
    return f(table, params)
  }

  // Push the message onto a queue to be processed serially.
  m := NewServerMessage(w, req, params, preprocess)
  s.channel <- m
  <- m.channel
}

// Processes a request within a table/servlet context.
func (s *Server) processWithObject(w http.ResponseWriter, req *http.Request, tableName string, objectId string, f func(*Table, *Servlet, map[string]interface{})(interface{}, error)) {
  params, err := decodeParams(w, req)
  if err != nil {
    return
  }

  // Create a wrapper function to open the table.
  preprocess := func(_ map[string]interface{})(interface{}, error) {
    table, err := s.OpenTable(tableName)
    if table == nil || err != nil {
      return nil, err
    }
    
    // Determine servlet index.
    index, err := s.GetObjectServletIndex(table, objectId)
    if err != nil {
      return nil, err
    }
    servlet := s.servlets[index]
    
    // Execute the original function.
    return f(table, servlet, params)
  }

  // Push the message onto a queue to be processed serially.
  m := NewServerMessage(w, req, params, preprocess)
  s.channel <- m
  <- m.channel
}

// Serially processes server messages routed through the server channel.
func (s *Server) processMessages() {
  for message := range s.channel {
    ret, err := message.f(message.params)
    if err != nil {
      s.writeResponse(message.w, message.req, map[string]interface{}{"message":err.Error()}, http.StatusInternalServerError)
    } else {
      s.writeResponse(message.w, message.req, ret, http.StatusOK)
    }
    message.channel <- true
  }
}

// Serially processes server messages routed through the server channel.
func (s *Server) writeResponse(w http.ResponseWriter, req *http.Request, ret interface{}, statusCode int) {
  if ret != nil {
    contentType := req.Header.Get("Content-Type")
    if contentType == "application/json" {
      encoder := json.NewEncoder(w)
      err := encoder.Encode(ret)
      if err != nil {
        fmt.Println(err.Error())
        return
      }
    }
  }
}

// Generates the path for a table attached to the server.
func (s *Server) GetTablePath(name string) string {
  return fmt.Sprintf("%v/%v", s.TablesPath(), name)
}

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
  table = NewTable(name, s.GetTablePath(name))
  err := table.Open()
  if err != nil {
    table.Close()
    return nil, err
  }
  s.tables[name] = table
  
  return table, nil
}


// Decodes the body of the message into parameters.
func decodeParams(w http.ResponseWriter, req *http.Request) (map[string]interface{}, error) {
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
