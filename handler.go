package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Maxtors/surisoc"
	"github.com/gorilla/mux"
)

// Response is the format we will use for all non socket responses
// the format is the same, so that any reciver will be able to handle them
type Response struct {
	Return  string `json:"return"`
	Message string `json:"message"`
}

func logger(inner http.HandlerFunc, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		NewLogItem(r, start).Log()
	})
}

func internalServerError(w http.ResponseWriter, message string) {
	bytes, err := json.MarshalIndent(&Response{Return: "NOK", Message: message}, "", "    ")
	if err != nil {
		log.Fatalf("Error when marshaling Response: %s\n", err.Error())
	}
	http.Error(w, string(bytes), http.StatusInternalServerError)
}

func handleCommand(w http.ResponseWriter, r *http.Request) {

	// Collect the command variable
	vars := mux.Vars(r)
	command := vars["command"]

	// Create a new Socket Message
	socketMessage := surisoc.NewSocketMessage(command)

	// Check if there are any arguments to parse
	if len(r.URL.Query()) != 0 {
		err := socketMessage.ParseArgumentsURLMap(r.URL.Query())
		if err != nil {
			internalServerError(w, fmt.Sprintf("Error: %s", err.Error()))
			return
		}
	}

	// Send the socket message
	response, err := session.SendMessage(socketMessage)
	if err != nil {
		internalServerError(w, fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	// Get the string representation of the response message
	res, err := response.ToString()
	if err != nil {
		internalServerError(w, fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	fmt.Fprintln(w, res)
}
