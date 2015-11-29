package main

import (
	"flag"
	"fmt"
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
	router.HandleFunc("/{command}/", handleCommand)

	// Start listening for requests
	log.Fatal(http.ListenAndServe(":8080", router))
}

func handleCommand(w http.ResponseWriter, r *http.Request) {

	// Collect the command variable
	vars := mux.Vars(r)
	command := vars["command"]

	// Create a new Socket Message
	socketMessage := surisoc.NewSocketMessage(command)

	// Check if there are any arguments to parse
	if len(r.URL.Query()) != 0 {
		err := socketMessage.ParseArgumentsMap(r.URL.Query())
		if err != nil {
			log.Fatalf("Error: %s", err.Error())
		}
	}

	// Send the socket message
	response, err := session.SendMessage(socketMessage)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}

	// Get the string representation of the response message
	res, err := response.ToString()
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
	fmt.Fprintln(w, res)
}
