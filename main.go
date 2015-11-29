package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Maxtors/surisoc"
	"github.com/gorilla/mux"
)

// Global variables for the Suricata Socket Application
var (
	socketPath string
	session    *surisoc.SuricataSocket
)

func init() {
	var err error

	// Parse commandline arguments
	flag.StringVar(&socketPath, "socket", "/var/run/suricata/suricata-command.socket", "Full path to the suricata unix socket")
	flag.Parse()

	// Create a new Suricata Socket session
	session, err = surisoc.NewSuricataSocket(socketPath)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}

func main() {
	defer session.Close()

	// Create a new router with one simple endpoint
	router := mux.NewRouter().StrictSlash(true)
	router.Handle("/{command}/", logger(handleCommand, "handle-command"))

	// Start listening for requests
	log.Fatal(http.ListenAndServe(":8080", router))
}
