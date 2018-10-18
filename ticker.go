package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO add t_stop field, whatever that means
// GET /api/ticker
func apiTickerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // The request has to be of GET type

		processStart := time.Now() // Track when the process started

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		timestamps := tickerTimestamps("")

		oldestTS := timestamps.oldestTimestamp
		latestTS := timestamps.latestTimestamp

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

		// t_stop SHOULD BE ADDED HERE

		response += `"tracks":` + `[`

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

		timestamps := tickerTimestamps("")
		latestTimestamp := timestamps.latestTimestamp

		if latestTimestamp.IsZero() { // If you dont assign a time to a time.Time variable, it's value is 0 date. We can check with IsZero() function
			fmt.Fprintln(w, "There are no track records")
		} else { //If it's not zero, we can format and display it to the user
			fmt.Fprintln(w, latestTimestamp.Format("02.01.2006 15:04:05.000"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

// GET /api/ticker/<timestamp>
func apiTickerTimestampHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // The request has to be of GET type

		processStart := time.Now() // Track when the process started

		pathArray := strings.Split(r.URL.Path, "/") // split the URL Path into chunks, whenever there's a "/"
		timestamp := pathArray[len(pathArray)-1]    // The part after the last "/", is the timestamp

		_, err := time.Parse("02.01.2006 15:04:05.000", timestamp) // Check if the timestamp provided is a valid time

		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // If there is an error, then return a bad request error
			return
		}

		timestamps := tickerTimestamps(timestamp)

		olderTS := timestamps.oldestNewerTimestamp
		latestTS := timestamps.latestTimestamp

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		response := `{`
		response += `"t_latest": "`
		if latestTS.IsZero() {
			response += `",`
		} else {
			response += latestTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		response += `"t_start": "`
		if olderTS.IsZero() {
			response += `",`
		} else {
			response += olderTS.Format("02.01.2006 15:04:05.000") + `",`
		}

		// t_stop SHOULD BE ADDED HERE

		response += `"tracks":` + `[`

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
