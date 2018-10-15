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

		// Check if track map contains the url
		// Or if the map is empty
		if !urlInMap(data["url"]) || len(igcFiles) == 0 {
			igcFiles[data["url"]] = URLTrack{"igc" + strconv.Itoa(igcFileCount), track, t} // Add the result to igcFiles map
			igcFileCount++                                                                 // Increase the count
		}

		response := `{`
		response += `"id": ` + `"` + igcFiles[data["url"]].trackName + `"`
		response += `}`

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON
		fmt.Fprintf(w, response)

	} else if r.Method == "GET" { // If the method is GET
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		x := 0 // Just some numeric iterator

		response := "["
		for i := range igcFiles { // Get all the IDs of .igc files stored in the igcFiles map
			if x != len(igcFiles)-1 { // If it's the last item in the array, don't add the ","
				response += "\"" + igcFiles[i].trackName + "\","
				x++ // Incerement the iterator
			} else {
				response += "\"" + igcFiles[i].trackName + "\""
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

		trackSliceURL := getTrackIndex(urlID)
		if trackSliceURL != "" { // Check whether the url is different from an empty string
			w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

			response := `{`
			response += `"H_date": "` + igcFiles[trackSliceURL].track.Date.String() + `",`
			response += `"pilot": "` + igcFiles[trackSliceURL].track.Pilot + `",`
			response += `"glider": "` + igcFiles[trackSliceURL].track.GliderType + `",`
			response += `"glider_id": "` + igcFiles[trackSliceURL].track.GliderID + `",`
			response += `"track_length": "` + calculateTotalDistance(igcFiles[trackSliceURL].track) + `",`
			response += `"track_src_url": "` + trackSliceURL + `"`
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

		trackSliceURL := getTrackIndex(uniqueID)

		if trackSliceURL != "" { // Check whether the url is different from an empty string

			something := map[string]string{ // Map the field to one of the Track struct attributes in the igcFiles slice
				"pilot":         igcFiles[trackSliceURL].track.Pilot,
				"glider":        igcFiles[trackSliceURL].track.GliderType,
				"glider_id":     igcFiles[trackSliceURL].track.GliderID,
				"track_length":  calculateTotalDistance(igcFiles[trackSliceURL].track),
				"H_date":        igcFiles[trackSliceURL].track.Date.String(),
				"track_src_url": trackSliceURL,
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
