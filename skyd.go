package skyd

import (
	"flag"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	defaultPort = 8585
	defaultDataDir = "/var/lib/sky"
)

const (
	portUsage = "the port to listen on"
	dataDirUsage = "the data directory"
)

//------------------------------------------------------------------------------
//
// Variables
//
//------------------------------------------------------------------------------

var port uint
var dataDir string

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

func init() {
	flag.UintVar(&port, "port", defaultPort, portUsage)
	flag.UintVar(&port, "p", defaultPort, portUsage+"(shorthand)")
	flag.StringVar(&dataDir, "data-dir", defaultDataDir, dataDirUsage)
	flag.StringVar(&dataDir, "d", defaultDataDir, dataDirUsage+"(shorthand)")
}

func main() {
	// Parse the command line arguments.
	flag.Parse()
	
	// Start the server up!
	c := make(chan bool)
	server := NewServer(port, dataDir)
	server.ListenAndServe(c)
	<- c
}