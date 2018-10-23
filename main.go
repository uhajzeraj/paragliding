package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var timeStarted = int(time.Now().Unix()) // Unix timestamp when the service started

// Gometalinter
const (
	gmlOB  = `{`
	gmlCB  = `}`
	gmlCPC = `",`
)

func urlRouter(w http.ResponseWriter, r *http.Request) {

	urlMap := map[string]func(http.ResponseWriter, *http.Request){ // A map of accepted URL RegEx patterns
		`^/paragliding$`:                          redirectHandler,
		`^/paragliding/api$`:                      apiHandler,
		`^/paragliding/api/track$`:                apiIgcHandler,
		`^/paragliding/api/track/igc[0-9]{1,10}$`: apiIgcIDHandler,
		`^/paragliding/api/track/igc[0-9]{1,10}/(pilot|glider|glider_id|track_length|H_date|track_src_url)$`: apiIgcIDFieldHandler,
		`^/paragliding/api/ticker/latest$`: apiTickerLatestHandler,
		`^/paragliding/api/ticker$`:        apiTickerHandler,
		`^/paragliding/api/ticker/\d{1,2}\.\d{1,2}\.\d{4} \d{1,2}:\d{1,2}:\d{1,2}.\d{1,3}$`: apiTickerTimestampHandler,
		`^/paragliding/api/webhook/new_track$`:                                              apiWebhookNewTrackHandler,
		`^/paragliding/api/webhook/new_track/webhook\d{1,3}$`:                               apiWebhookNewTrackWebhookIDHandler,
		`^/paragliding/admin/api/tracks_count$`:                                             adminAPITrackCountHandler,
		`^/paragliding/admin/api/tracks$`:                                                   adminAPITracksDelete,

		// This is a special case used to trigger the webhooks
		`^/paragliding/admin/api/webhook$`: adminAPIWebhookTrigger,
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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
