package main

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
)

// Timestamps for ticker API struct
type Timestamps struct {
	latestTimestamp      time.Time
	oldestTimestamp      time.Time
	oldestNewerTimestamp time.Time
}

// ISO8601 duration parsing function
func parseTimeDifference(timeDifference int) string {

	result := "P" // Different time intervals are attached to this, if they are != 0

	// Lower the cyclomatic complexity
	// Save the durations and keys in arrays

	var dateValue [4]int
	dateKey := [4]string{"Y", "M", "W", "D"}

	var timeValue [3]int
	timeKey := [3]string{"H", "M", "S"}

	// Formulas for calculating different time intervals in seconds

	// Date
	timeLeft := timeDifference
	dateValue[0] = timeDifference / 31557600

	timeLeft -= dateValue[0] * 31557600
	dateValue[1] = timeLeft / 2592000

	timeLeft -= dateValue[1] * 2592000
	dateValue[2] = timeLeft / 604800

	timeLeft -= dateValue[2] * 604800
	dateValue[3] = timeLeft / 86400

	// Time
	timeLeft -= dateValue[3] * 86400
	timeValue[0] = timeLeft / 3600

	timeLeft -= timeValue[0] * 3600
	timeValue[1] = timeLeft / 60

	timeLeft -= timeValue[1] * 60
	timeValue[2] = timeLeft

	for i := 0; i < 4; i++ {
		// Add the time intervals if they are diffrent from 0
		if dateValue[i] != 0 {
			result += dateKey[i] + strconv.Itoa(dateValue[i])
		}
	}

	// // Check in case time intervals are 0
	if timeValue[0] != 0 || timeValue[1] != 0 || timeValue[2] != 0 {
		result += "T"
	}

	for i := 0; i < 3; i++ {
		// Add the time intervals if they are diffrent from 0
		if timeValue[i] != 0 {
			result += timeKey[i] + strconv.Itoa(timeValue[i])
		}
	}

	return result

}

// TODO add only 2 middle points (3 distances in total) **OPTIONAL?**
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
func latestTimestamp(resultTracks []igcTrack) time.Time {
	var latestTimestamp time.Time // Create a variable to store the most recent track added

	for _, val := range resultTracks { // Iterate every track to find the most recent track added
		if val.TimeRecorded.After(latestTimestamp) { // If current track timestamp is after the current latestTimestamp...
			latestTimestamp = val.TimeRecorded // Set that one as the latestTimestamp
		}
	}

	return latestTimestamp
}

// Return the oldest timestamp
func oldestTimestamp(resultTracks []igcTrack) time.Time {

	// Just the first time, add the first found timestamp
	// After that, check that one against the other timestamps in the slice
	// If there is none, JSON response will be an empty string ""
	// If there is one timestamp, that one is the oldest timestamp as well

	var oldestTimestamp time.Time // Create a variable to store the oldest track added

	for key, val := range resultTracks { // Iterate every track to find the oldest track added

		// Assign to oldestTimestamp a value, but just once
		// Then we check it against other timestamps of other tracks in the slice
		if key == 0 {
			oldestTimestamp = val.TimeRecorded
		}

		if val.TimeRecorded.Before(oldestTimestamp) { // If current track timestamp is before the current latestTimestamp...
			oldestTimestamp = val.TimeRecorded // Set that one as the latestTimestamp
		}
	}

	return oldestTimestamp
}

// Return the oldest timestamp which is newer than input timestamp
func oldestNewerTimestamp(inputTS string, resultTracks []igcTrack) time.Time {

	ts := time.Now()
	testTs := ts

	parsedTime, _ := time.Parse("02.01.2006 15:04:05.000", inputTS) // Parse the string into time

	for _, val := range resultTracks { // Iterate every track to find the most recent track added
		if val.TimeRecorded.After(parsedTime) && val.TimeRecorded.Before(ts) { // If current track timestamp is after the current latestTimestamp...
			ts = val.TimeRecorded // Set that one as the latestTimestamp
		}
	}

	if testTs.Equal(ts) {
		return time.Time{}
	}

	return ts
}

func tickerTimestamps(inputTS string) Timestamps {
	conn := mongoConnect()
	resultTracks := getAllTracks(conn, false)

	timestamps := Timestamps{}

	timestamps.latestTimestamp = latestTimestamp(resultTracks)
	timestamps.oldestTimestamp = oldestTimestamp(resultTracks)
	timestamps.oldestNewerTimestamp = oldestNewerTimestamp(inputTS, resultTracks)

	return timestamps
}

// Return track names
// And also t_stop track
func returnTracks(n int) (string, time.Time) {
	var response string
	var tStop time.Time

	conn := mongoConnect()

	resultTracks := getAllTracks(conn, false)

	for key, val := range resultTracks { // Go through the slice
		response += `"` + val.TrackName + `",`
		if key == n-1 || key == len(resultTracks)-1 {
			tStop = val.TimeRecorded
			break
		}
	}

	// Get rid of that last `,` of JSON will freak out
	response = strings.TrimRight(response, ",")

	return response, tStop
}
