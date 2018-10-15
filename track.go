package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
)

// URLTrack - Keep track of the url used for adding the igc file
type igcTrack struct {
	trackURL        string
	trackName       string
	timeRecorded    time.Time
	trackDate       time.Time
	trackPilot      string
	trackGliderType string
	trackGliderID   string
	trackPoints     []igc.Point
}

var igcTrackCount = 1 // Keep count of the number of igc files added to the system

// Map where the igcFiles are in-memory stored
var igcTracks []igcTrack // slice of igcTrack

//
//
//
// Track Handlers

// POST/GET /api/track
func apiIgcHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" { // If method is POST, user has entered the URL
		var data map[string]string // POST body is of content-type: JSON; the result can be stored in a map
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			return
		}

		track, err := igc.ParseLocation(data["url"]) // call the igc library
		if err != nil {
			return
		}

		t := time.Now()

		data["url"] = strings.Replace(data["url"], "%20", " ", -1) // Replace %20 with " " in the URL

		// Check if track slice contains the url
		// Or if the slice is empty
		if !urlInSlice(data["url"]) || len(igcTracks) == 0 {
			igcTracks = append(igcTracks, // Append the result to igcTracks slice
				igcTrack{data["url"], "igc" + strconv.Itoa(igcTrackCount), t, track.Date, track.Pilot, track.GliderType, track.GliderID, track.Points})
			igcTrackCount++ // Increase the count
		}

		// Find the id of the track in igcTracks slice
		trackID := idOfTrack(data["url"])

		response := `{`
		response += `"id": ` + `"` + igcTracks[trackID].trackName + `"`
		response += `}`

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON
		fmt.Fprintf(w, response)

	} else if r.Method == "GET" { // If the method is GET
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		x := 0 // Just some numeric iterator

		response := "["
		for i := range igcTracks { // Get all the IDs of .igc files stored in the igcFiles map
			if x != len(igcTracks)-1 { // If it's the last item in the array, don't add the ","
				response += "\"" + igcTracks[i].trackName + "\","
				x++ // Incerement the iterator
			} else {
				response += "\"" + igcTracks[i].trackName + "\""
			}
		}
		response += "]"

		fmt.Fprintf(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't any of those, send a 404 Not Found status
	}
}

// GET /api/track/<igcFile>
func apiIgcIDHandler(w http.ResponseWriter, r *http.Request) {

	// The request has to be of GET type
	if r.Method == "GET" {
		urlID := path.Base(r.URL.Path) // returns the part after the last '/' in the url

		igcTracksIndex := getTrackIndex(urlID)
		if igcTracksIndex != -1 { // Check whether the url is different from an empty string
			w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

			response := `{`
			response += `"H_date": "` + igcTracks[igcTracksIndex].trackDate.String() + `",`
			response += `"pilot": "` + igcTracks[igcTracksIndex].trackPilot + `",`
			response += `"glider": "` + igcTracks[igcTracksIndex].trackGliderType + `",`
			response += `"glider_id": "` + igcTracks[igcTracksIndex].trackGliderID + `",`
			response += `"track_length": "` + calculateTotalDistance(igcTracks[igcTracksIndex].trackPoints) + `",`
			response += `"track_src_url": "` + igcTracks[igcTracksIndex].trackURL + `"`
			response += `}`

			fmt.Fprintf(w, response)
		} else {
			w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

// GET /api/track/<igcFile>/<field>
func apiIgcIDFieldHandler(w http.ResponseWriter, r *http.Request) {

	// The request has to be of GET type
	if r.Method == "GET" {
		pathArray := strings.Split(r.URL.Path, "/") // split the URL Path into chunks, whenever there's a "/"
		field := pathArray[len(pathArray)-1]        // The part after the last "/", is the field
		uniqueID := pathArray[len(pathArray)-2]     // The part after the second to last "/", is the unique ID

		igcTracksIndex := getTrackIndex(uniqueID)

		if igcTracksIndex != -1 { // Check whether the url is different from an empty string

			something := map[string]string{ // Map the field to one of the Track struct attributes in the igcFiles slice
				"pilot":         igcTracks[igcTracksIndex].trackPilot,
				"glider":        igcTracks[igcTracksIndex].trackGliderType,
				"glider_id":     igcTracks[igcTracksIndex].trackGliderID,
				"track_length":  calculateTotalDistance(igcTracks[igcTracksIndex].trackPoints),
				"H_date":        igcTracks[igcTracksIndex].trackDate.String(),
				"track_src_url": igcTracks[igcTracksIndex].trackURL,
			}

			response := something[field] // This will work because the RegEx checks whether the name is written correctly
			fmt.Fprintf(w, response)
		} else {
			w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}
