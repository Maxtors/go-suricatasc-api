package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Maxtors/surisoc"
	"github.com/gorilla/mux"
)

// Global variables for the Suricata Socket Application
var (
	socketPath     string
	session        *surisoc.SuricataSocket
	bindingAddress string
	signals        chan os.Signal
	done           chan bool
)

func init() {
	log.Println("Welcome to go-suricatasc-api!")
	var err error

	// Parse commandline arguments
	flag.StringVar(&socketPath, "socket", "/var/run/suricata/suricata-command.socket", "Full path to the suricata unix socket")
	host := flag.String("host", "127.0.0.1", "The IP-Address to bind to")
	port := flag.String("port", "8080", "The Port to bind to")
	flag.Parse()

	if host != nil && port != nil {
		bindingAddress = net.JoinHostPort(*host, *port)
	}

	// Create a new Suricata Socket session
	session, err = surisoc.NewSuricataSocket(socketPath)
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}

	// Create channels
	signals = make(chan os.Signal, 1)
	done = make(chan bool, 1)

	// Set up the signals to listen to, SIGQUIT is used internaly to stop
	// the current process
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	log.Println("Done initializing")
}

func main() {
	defer session.Close()

	// Go-Routine to handle reciving of signals
	go func() {
		signal := <-signals
		log.Printf("Recived signal: %s\n", signal)
		done <- true
	}()

	// Go-Routine to handle the API functionality
	go func() {
		// Create a new router with one simple endpoint
		router := mux.NewRouter().StrictSlash(true)
		router.Handle("/{command}/", logger(handleCommand, "handle-command"))

		// Start listening for requests
		log.Printf("Started listening to requests [%s]", bindingAddress)
		log.Println(http.ListenAndServe(bindingAddress, router))
		signals <- syscall.SIGQUIT
	}()

	// Wait for a done signal
	<-done
	log.Println("Goodbye!")
}
