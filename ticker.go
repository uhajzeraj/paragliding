package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// TODO add t_stop field, whatever that means
// GET /api/ticker
func apiTickerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // The request has to be of GET type

		processStart := time.Now() // Track when the process started

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		oldestTS := oldestTimestamp()
		latestTS := latestTimestamp()

		// timestamps := returnTimestamps(5)

		response := `{`
		response += `"t_latest": "`
		if latestTS.IsZero() {
			response += `",`
		} else {
			response += latestTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		response += `"t_start": "`
		if oldestTS.IsZero() {
			response += `",`
		} else {
			response += oldestTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		response += `"tracks":` + `[`
		// t_stop SHOULD BE ADDED HERE

		// THAT 5 SHOULD BE ABLE TO CHANGE DYNAMICALLY
		response += returnTracks(5) // Maximum of 5 tracks

		response += `],`
		response += `"processing":` + `"` + strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 6, 64) + `ms"`
		response += `}`
		fmt.Fprintln(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}

}

// GET /api/ticker/latest
func apiTickerLatestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // The request has to be of GET type
		latestTimestamp := latestTimestamp()

		if latestTimestamp.IsZero() { // If you dont assign a time to a time.Time variable, it's value is 0 date. We can check with IsZero() function
			fmt.Fprintln(w, "There are no track records")
		} else { //If it's not zero, we can format and display it to the user
			fmt.Fprintln(w, latestTimestamp.Format("02.01.2006 15:04:05.000"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}
