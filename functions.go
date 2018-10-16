package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	igc "github.com/marni/goigc"
)

// Find the id of the track against URL in the igcTracks slice
func idOfTrack(url string) int {
	for key, val := range igcTracks {
		if val.trackURL == url {
			return key
		}
	}
	return -1
}

// Check if url is in the igcTracks slice
func urlInSlice(url string) bool {
	for _, urlInMap := range igcTracks {
		if urlInMap.trackURL == url {
			return true
		}
	}
	return false
}

// Get the index of the track in the igcFiles slice, if it is there
func getTrackIndex(trackName string) int {
	for key, track := range igcTracks {
		if track.trackName == trackName {
			return key
		}
	}
	return -1
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
func calculateTotalDistance(trackPoints []igc.Point) string {

	totalDistance := 0.0

	// For each point of the track, calculate the distance between 2 points in the Point array
	for i := 0; i < len(trackPoints)-1; i++ {
		totalDistance += trackPoints[i].Distance(trackPoints[i+1])
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

// Return the latest timestamp
func latestTimestamp() time.Time {
	var latestTimestamp time.Time   // Create a variable to store the most recent track added
	for _, val := range igcTracks { // Iterate every track to find the most recent track added
		if val.timeRecorded.After(latestTimestamp) { // If current track timestamp is after the current latestTimestamp...
			latestTimestamp = val.timeRecorded // Set that one as the latestTimestamp
		}
	}

	return latestTimestamp
}

// Return the oldest timestamp
func oldestTimestamp() time.Time {

	// Just the first time, add the first found timestamp
	// After that, check that one against the other timestamps in the slice
	// If there is none, JSON response will be an empty string ""
	// If there is one timestamp, that one is the oldest timestamp as well

	var oldestTimestamp time.Time     // Create a variable to store the oldest track added
	for key, val := range igcTracks { // Iterate every track to find the oldest track added

		// Assign to oldestTimestamp a value, but just once
		// Then we check it against other timestamps of other tracks in the slice
		if key == 0 {
			oldestTimestamp = val.timeRecorded
		}

		if val.timeRecorded.Before(oldestTimestamp) { // If current track timestamp is before the current latestTimestamp...
			oldestTimestamp = val.timeRecorded // Set that one as the latestTimestamp
		}
	}

	return oldestTimestamp
}

// Return the oldest timestamp which is newer than input timestamp
func oldestNewerTimestamp(inputTS string) time.Time {

	ts := time.Now()
	testTs := ts

	parsedTime, _ := time.Parse("02.01.2006 15:04:05.000", inputTS) // Parse the string into time

	for _, val := range igcTracks { // Iterate every track to find the most recent track added
		if val.timeRecorded.After(parsedTime) && val.timeRecorded.Before(ts) { // If current track timestamp is after the current latestTimestamp...
			ts = val.timeRecorded // Set that one as the latestTimestamp
		}
	}

	if testTs.Equal(ts) {
		return time.Time{}
	}

	return ts
}

// Return track names
func returnTracks(n int) string {
	response := ""

	for key, val := range igcTracks { // Go through the slice
		if key < n-1 { // Check if the count is less than the number of required elements
			if key == len(igcTracks)-1 {
				response += `"` + val.trackName + `"` // Append the tackName to the response
				break                                 // Break out of the loop, no need to add any other elements
			} else {
				response += `"` + val.trackName + `",` // Append the trackName to the response
			}
		} else {
			response += `"` + val.trackName + `"` // Append the tackName to the response
			break                                 // Break out of the loop, no need to add any other elements
		}
	}

	return response
}
