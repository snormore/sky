package main

import (
	"flag"
	"fmt"
	"github.com/skydb/sky/skyd"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
)

//------------------------------------------------------------------------------
//
// Variables
//
//------------------------------------------------------------------------------

var config *skyd.Config
var configPath string

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

//--------------------------------------
// Initialization
//--------------------------------------

func init() {
	config = skyd.NewConfig()
	flag.UintVar(&config.Port, "port", config.Port, "the port to listen on")
	flag.UintVar(&config.Port, "p", config.Port, "the port to listen on")
	flag.StringVar(&config.DataPath, "data-path", config.DataPath, "the data directory")
	flag.StringVar(&config.PidPath, "pid-path", config.PidPath, "the path to the pid file")
	flag.StringVar(&configPath, "config", "", "the path to the config file")
}

//--------------------------------------
// Main
//--------------------------------------

func main() {
	// Parse the command line arguments and load the config file (if specified).
	flag.Parse()
	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Printf("Unable to open config: %v\n", err)
			return
		}
		defer file.Close()
		if err = config.Decode(file); err != nil {
			fmt.Printf("Unable to parse config: %v\n", err)
			os.Exit(1)
		}
	}

	// Hardcore parallelism right here.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize
	server := skyd.NewServer(config.Port, config.DataPath)
	writePidFile()
	setupSignalHandlers(server)

	// Start the server up!
	c := make(chan bool)
	err := server.ListenAndServe(c)
	if err != nil {
		fmt.Printf("%v\n", err)
		cleanup(server)
		return
	}
	<-c
	cleanup(server)
}

//--------------------------------------
// Signals
//--------------------------------------

// Handles signals received from the OS.
func setupSignalHandlers(server *skyd.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			fmt.Fprintln(os.Stderr, "Shutting down...")
			cleanup(server)
			fmt.Fprintln(os.Stderr, "Shutdown complete.")
			os.Exit(1)
		}
	}()
}

//--------------------------------------
// Utility
//--------------------------------------

// Shuts down the server socket and closes the database.
func cleanup(server *skyd.Server) {
	if server != nil {
		server.Shutdown()
	}
	deletePidFile()
}

// Writes a file to /var/run that contains the current process id.
func writePidFile() {
	pid := fmt.Sprintf("%d", os.Getpid())
	if err := ioutil.WriteFile(config.PidPath, []byte(pid), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write pid file: %v\n", err)
	}
}

// Deletes the pid file.
func deletePidFile() {
	if _, err := os.Stat(config.PidPath); !os.IsNotExist(err) {
		if err = os.Remove(config.PidPath); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to remove pid file: %v\n", err)
		}
	}
}
