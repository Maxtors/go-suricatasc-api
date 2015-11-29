package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// LogItem is an event in the API that will be logged
type LogItem struct {
	Timestamp  time.Time     `json:"timestamp"`
	Method     string        `json:"http_method"`
	URL        string        `json:"http_url"`
	Duration   time.Duration `json:"duration"`
	RemoteAddr string        `json:"http_remote_address"`
	Protocol   string        `json:"http_proto"`
	UserAgent  string        `json:"user_agent"`
}

// NewLogItem creates a new LogItem based on a http request
func NewLogItem(r *http.Request, startTime time.Time) *LogItem {
	l := &LogItem{
		Timestamp:  time.Now(),
		Method:     r.Method,
		URL:        r.RequestURI,
		RemoteAddr: r.RemoteAddr,
		Duration:   time.Since(startTime),
		Protocol:   r.Proto,
		UserAgent:  r.UserAgent(),
	}
	return l
}

// Log will log a LogItem to the logger
func (l LogItem) Log() {
	bytes, err := json.Marshal(l)
	if err != nil {
		log.Fatalf("Error when marshaling LogItem: %s\n", err.Error())
	}

	fmt.Printf("%s\n", string(bytes))
}
