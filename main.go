package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
)

// URLTrack - Keep track of the url used for adding the igc file
type URLTrack struct {
	trackName    string
	track        igc.Track
	timeRecorded time.Time
}

// Keep count of the number of igc files added to the system
var igcFileCount = 1

// Map where the igcFiles are in-memory stored
var igcFiles = make(map[string]URLTrack) // map["URL"]urlTrack

// Unix timestamp when the service started
var timeStarted = int(time.Now().Unix())

// Check if url is in the urlTrack map
func urlInMap(url string) bool {
	for urlInMap := range igcFiles {
		if urlInMap == url {
			return true
		}
	}
	return false
}

// Get the index of the track in the igcFiles slice, if it is there
func getTrackIndex(trackName string) string {
	for url, track := range igcFiles {
		if track.trackName == trackName {
			return url
		}
	}
	return ""
}

// ISO8601 duration parsing function
func parseTimeDifference(timeDifference int) string {

	result := "P" // Different time intervals are attached to this, if they are != 0

	// Formulas for calculating different time intervals in seconds
	timeLeft := timeDifference
	years := timeDifference / 31557600
	timeLeft -= years * 31557600
	months := timeLeft / 2592000
	timeLeft -= months * 2592000
	weeks := timeLeft / 604800
	timeLeft -= weeks * 604800
	days := timeLeft / 86400
	timeLeft -= days * 86400
	hours := timeLeft / 3600
	timeLeft -= hours * 3600
	minutes := timeLeft / 60
	timeLeft -= minutes * 60
	seconds := timeLeft

	// Add time invervals to the result only if they are different form 0
	if years != 0 {
		result += fmt.Sprintf("Y%d", years)
	}
	if months != 0 {
		result += fmt.Sprintf("M%d", months)
	}
	if weeks != 0 {
		result += fmt.Sprintf("W%d", weeks)
	}
	if days != 0 {
		result += fmt.Sprintf("D%d", days)
	}

	if hours != 0 || minutes != 0 || seconds != 0 { // Check in case time intervals are 0
		result += "T"
		if hours != 0 {
			result += fmt.Sprintf("H%d", hours)
		}
		if minutes != 0 {
			result += fmt.Sprintf("M%d", minutes)
		}
		if seconds != 0 {
			result += fmt.Sprintf("S%d", seconds)
		}
	}

	return result
}

// Calculate the total distance of the track
func calculateTotalDistance(track igc.Track) string {

	totalDistance := 0.0

	// For each point of the track, calculate the distance between 2 points in the Point array
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	// Parse it to a string value
	return strconv.FormatFloat(totalDistance, 'f', 2, 64)
}

// Check if any of the regex patterns supplied in the map parameter match the string parameter
func regexMatches(url string, urlMap map[string]func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	for mapURL := range urlMap {
		res, err := regexp.MatchString(mapURL, url)
		if err != nil {
			return nil
		}

		if res { // If the pattern matching returns true, return the function
			return urlMap[mapURL]
		}
	}
	return nil
}

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

		response := "{"
		response += "\"id\": " + "\"" + igcFiles[data["url"]].trackName + "\""
		response += "}"

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON
		fmt.Fprintf(w, response)

		fmt.Println("Time: ", igcFiles[data["url"]].timeRecorded.Format("02.01.2006 15:04:05.000"))

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

func apiIgcIDHandler(w http.ResponseWriter, r *http.Request) {

	// The request has to be of GET type
	if r.Method == "GET" {
		urlID := path.Base(r.URL.Path) // returns the part after the last '/' in the url

		trackSliceURL := getTrackIndex(urlID)
		if trackSliceURL != "" { // Check whether the url is different from an empty string
			w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

			response := "{"
			response += "\"H_date\": " + "\"" + igcFiles[trackSliceURL].track.Date.String() + "\","
			response += "\"pilot\": " + "\"" + igcFiles[trackSliceURL].track.Pilot + "\","
			response += "\"glider\": " + "\"" + igcFiles[trackSliceURL].track.GliderType + "\","
			response += "\"glider_id\": " + "\"" + igcFiles[trackSliceURL].track.GliderID + "\","
			response += "\"track_length\": " + "\"" + calculateTotalDistance(igcFiles[trackSliceURL].track) + "\","
			response += "\"track_src_url\": " + "\"" + trackSliceURL + "\""
			response += "}"

			fmt.Fprintf(w, response)
		} else {
			w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

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

// GET /api
func apiHandler(w http.ResponseWriter, r *http.Request) {

	// The request has to be of GET type
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		timeNow := int(time.Now().Unix()) // Unix timestamp when the handler was called

		iso8601duration := parseTimeDifference(timeNow - timeStarted) // Calculate the time elapsed by subtracting the times

		response := "{"
		response += "\"uptime\": \"" + iso8601duration + "\","
		response += "\"info\": \"Service for IGC tracks.\","
		response += "\"version\": \"v1\""
		response += "}"
		fmt.Fprintln(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

// Redirect to /api
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "paragliding/api", http.StatusSeeOther) // Redirect this request
}

// GET /api/ticker/latest
func apiTickerLatestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { // The request has to be of GET type
		var latestTimestamp time.Time  // Create a variable to store the most recent track added
		for _, val := range igcFiles { // Iterate every track to find the most recent track added
			if val.timeRecorded.After(latestTimestamp) { // If current track timestamp is after the current latestTimestamp...
				latestTimestamp = val.timeRecorded // Set that one as the latestTimestamp
			}
		}

		if latestTimestamp.IsZero() { // If you dont assign a time to a time.Time variable, it's value is 0 date. We can check with IsZero() function
			fmt.Fprintln(w, "There are no track records")
		} else { //If it's not zero, we can format and display it to the user
			fmt.Fprintln(w, latestTimestamp.Format("02.01.2006 15:04:05.000"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

func urlRouter(w http.ResponseWriter, r *http.Request) {

	urlMap := map[string]func(http.ResponseWriter, *http.Request){ // A map of accepted URL RegEx patterns
		"^/paragliding$":                          redirectHandler,
		"^/paragliding/api$":                      apiHandler,
		"^/paragliding/api/track$":                apiIgcHandler,
		"^/paragliding/api/track/igc[0-9]{1,10}$": apiIgcIDHandler,
		"^/paragliding/api/track/igc[0-9]{1,10}/(pilot|glider|glider_id|track_length|H_date|track_src_url)$": apiIgcIDFieldHandler,
		"^/paragliding/api/ticker/latest$": apiTickerLatestHandler,
		// "^/"
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
