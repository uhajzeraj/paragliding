package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// // Webhook - Handle the decoded JSON body from the POST request
// type Webhook struct {
// 	WebhookURL      string `json:"webhook_url"`
// 	MinTriggerValue int    `json:"min_trigger_value"`
// 	WebhookID       string `json:"webhook_id"`
// }

// POST /api/webhook/new_track
func apiWebhookNewTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // The method has to be POST
		var data map[string]interface{} // Save POST body (which is JSON) in the map

		err := json.NewDecoder(r.Body).Decode(&data) // Json decode
		if err != nil {
			return
		}

		// Check if webhookURL exists in the map
		// If it doesn't, return
		if _, ok := data["webhookURL"]; !ok {
			return
		}

		// Check if minTriggerValue exists in the map
		// If there is not, set the minTriggerValue to 1
		if _, ok := data["minTriggerValue"]; !ok {
			data["minTriggerValue"] = 1
		}

		webhookID := insertUpdateWebhook(data)

		// Print the webhookID
		fmt.Fprintln(w, webhookID)

	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}

// DELETE/GET /api/webhook/new_track/<webhook_id>
func apiWebhookNewTrackWebhookIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet || r.Method == http.MethodDelete { // GET or DELETE method

		// The methods are almost the same
		// DELETE has to additionally delete the webhook

		pathArray := strings.Split(r.URL.Path, "/") // split the URL Path into chunks, whenever there's a "/"
		webhookID := pathArray[len(pathArray)-1]    // The part after the last "/", is the webhook_id

		webhook := getWebhook(mongoConnect(), webhookID)

		if webhook.WebhookID != "" { // Check whether the name is different from an empty string

			w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

			response := gmlOB
			response += `"webhookURL": "` + webhook.WebhookURL + `",`
			response += `"minTriggerValue": ` + strconv.Itoa(webhook.MinTriggerValue)
			response += gmlCB

			fmt.Fprintln(w, response)
		} else {
			w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
		}

		// If the method was DELETE, delete the webhook
		if r.Method == http.MethodDelete {
			deleteWebhook(mongoConnect(), webhookID)
		}

	} else { // Method is something different
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}

// Trigger the webhook(s) if new tracks were added
func triggerWebhook() {

	conn := mongoConnect()
	resultWebhooks := getAllWebhooks(conn)
	trackCounter := getCounter(conn.Database("paragliding"), "trackCounter") - 1

	for _, val := range resultWebhooks {

		// Only check for new tracks
		if trackCounter%val.MinTriggerValue != 0 {
			continue
		}

		processStart := time.Now() // Track when the process started

		url := val.WebhookURL

		// returnTracks returns the last element and the n number of tracks
		trackString, _ := returnTracks(trackCounter)

		trackArray := strings.Split(trackString, `,`)

		if len(trackArray) < val.MinTriggerValue {
			trackArray = trackArray[0:len(trackArray)]
		} else {
			trackArray = trackArray[trackCounter-val.MinTriggerValue : len(trackArray)]
		}

		trackString = strings.Join(trackArray, `,`)
		// Add \ to " in trackArray
		trackString = strings.Replace(trackString, `"`, `\"`, -1)

		timestamps := tickerTimestamps("")

		latestTS := timestamps.latestTimestamp.String()
		jsonPayload := gmlOB
		jsonPayload += `"username": "Track Added",`
		jsonPayload += `"content": "` + latestTS + `\n`
		jsonPayload += `[` + trackString + `]\n`
		jsonPayload += strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + `ms"`
		jsonPayload += gmlCB

		var jsonStr = []byte(jsonPayload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}
