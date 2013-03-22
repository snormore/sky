package main

import (
	"./skyd"
	"flag"
	"runtime"
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
	
	// Hardcore parallelism right here.
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	// Start the server up!
	c := make(chan bool)
	server := skyd.NewServer(port, dataDir)
	server.ListenAndServe(c)
	<- c
}