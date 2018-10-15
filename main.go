package main

import (
	"log"
	"net/http"
	"time"

	igc "github.com/marni/goigc"
)

// URLTrack - Keep track of the url used for adding the igc file
type URLTrack struct {
	trackName    string
	track        igc.Track
	timeRecorded time.Time
}

var igcFileCount = 1 // Keep count of the number of igc files added to the system

// Map where the igcFiles are in-memory stored
var igcFiles = make(map[string]URLTrack) // map["URL"]urlTrack

var timeStarted = int(time.Now().Unix()) // Unix timestamp when the service started

func urlRouter(w http.ResponseWriter, r *http.Request) {

	urlMap := map[string]func(http.ResponseWriter, *http.Request){ // A map of accepted URL RegEx patterns
		"^/paragliding$":                          redirectHandler,
		"^/paragliding/api$":                      apiHandler,
		"^/paragliding/api/track$":                apiIgcHandler,
		"^/paragliding/api/track/igc[0-9]{1,10}$": apiIgcIDHandler,
		"^/paragliding/api/track/igc[0-9]{1,10}/(pilot|glider|glider_id|track_length|H_date|track_src_url)$": apiIgcIDFieldHandler,
		"^/paragliding/api/ticker/latest$": apiTickerLatestHandler,
		"^/paragliding/api/ticker$":        apiTickerHandler,
	}

	result := regexMatches(r.URL.Path, urlMap) // Perform the RegEx check to see if any pattern matches

	if result != nil { // If a function is returned, call that handler function
		result(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

func main() {
	http.HandleFunc("/", urlRouter) // Handle all the request via the urlRouter function
	log.Fatal(http.ListenAndServe(":8080", nil))
}
