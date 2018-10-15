package main

import (
	"fmt"
	"net/http"
	"time"
)

// GET /api
func apiHandler(w http.ResponseWriter, r *http.Request) {

	// The request has to be of GET type
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		timeNow := int(time.Now().Unix()) // Unix timestamp when the handler was called

		iso8601duration := parseTimeDifference(timeNow - timeStarted) // Calculate the time elapsed by subtracting the times

		response := `{`
		response += `"uptime": "` + iso8601duration + `",`
		response += `"info": "Service for IGC tracks.",`
		response += `"version": "v1"`
		response += `}`
		fmt.Fprintln(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

// Redirect to /api
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "paragliding/api", http.StatusSeeOther) // Redirect this request
}
