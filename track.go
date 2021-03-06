package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
)

// URLTrack - Keep track of the url used for adding the igc file
type igcTrack struct {
	TrackURL           string
	TrackName          string
	TimeRecorded       time.Time
	TrackDate          time.Time
	TrackPilot         string
	TrackGliderType    string
	TrackGliderID      string
	TrackTotalDistance string
}

// POST/GET /api/track
func apiIgcHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost { // If method is POST, user has entered the URL
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

		// Connect to MongoDB
		conn := mongoConnect()

		// Choose database and collection
		db := conn.Database("paragliding")  // `paragliding` database
		trackColl := db.Collection("track") // `track` collection

		id := getCounter(db, "trackCounter")

		// Check if track is already in the database
		if !urlInMongo(data["url"], trackColl) { // If it is not, we can add it

			totDist := calculateTotalDistance(track.Points)

			// The info which needs to be saved in the database
			trackWrite := igcTrack{TrackURL: data["url"], TrackName: "igc" + strconv.Itoa(id), TimeRecorded: t,
				TrackDate: track.Date, TrackPilot: track.Pilot, TrackGliderType: track.GliderType, TrackGliderID: track.GliderID, TrackTotalDistance: totDist}
			// Insert it in the database
			_, err = trackColl.InsertOne(context.Background(), trackWrite)
			if err != nil {
				log.Fatal(err)
				return
			}

			// Increase the counter stored in the database
			increaseCounter(int32(id), db, "trackCounter")

			// Trigger the webhook
			triggerWebhook()
		}

		resultTrackName := trackNameFromURL(data["url"], trackColl)

		response := gmlOB
		response += `"id": "` + resultTrackName + `"`
		response += gmlCB

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON
		fmt.Fprintf(w, response)

	} else if r.Method == http.MethodGet { // If the method is GET
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		x := 0 // Just some numeric iterator

		// Connect to MongoDB
		conn := mongoConnect()

		igcTracks := getAllTracks(conn)

		response := "["
		for i := range igcTracks { // Get all the IDs of .igc files stored in the igcFiles map
			if x != len(igcTracks)-1 { // If it's the last item in the array, don't add the ","
				response += "\"" + igcTracks[i].TrackName + "\","
				x++ // Incerement the iterator
			} else {
				response += "\"" + igcTracks[i].TrackName + "\""
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
	if r.Method == http.MethodGet {
		urlTrackName := path.Base(r.URL.Path) // returns the part after the last '/' in the url

		conn := mongoConnect()

		resIgcTrack := getTrack(conn, urlTrackName)

		if resIgcTrack.TrackName != "" { // Check whether the name is different from an empty string
			w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

			response := gmlOB
			response += `"H_date": "` + resIgcTrack.TrackDate.String() + `",`
			response += `"pilot": "` + resIgcTrack.TrackPilot + `",`
			response += `"glider": "` + resIgcTrack.TrackGliderType + `",`
			response += `"glider_id": "` + resIgcTrack.TrackGliderID + `",`
			response += `"track_length": "` + resIgcTrack.TrackTotalDistance + `",`
			response += `"track_src_url": "` + resIgcTrack.TrackURL + `"`
			response += gmlCB

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
	if r.Method == http.MethodGet {
		pathArray := strings.Split(r.URL.Path, "/") // split the URL Path into chunks, whenever there's a "/"
		field := pathArray[len(pathArray)-1]        // The part after the last "/", is the field
		uniqueID := pathArray[len(pathArray)-2]     // The part after the second to last "/", is the unique ID

		conn := mongoConnect()

		resIgcTrack := getTrack(conn, uniqueID)

		if resIgcTrack.TrackName != "" { // Check whether the name is different from an empty string

			something := map[string]string{ // Map the field to one of the Track struct attributes in the igcFiles slice
				"pilot":         resIgcTrack.TrackPilot,
				"glider":        resIgcTrack.TrackGliderType,
				"glider_id":     resIgcTrack.TrackGliderID,
				"track_length":  resIgcTrack.TrackTotalDistance,
				"H_date":        resIgcTrack.TrackDate.String(),
				"track_src_url": resIgcTrack.TrackURL,
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
