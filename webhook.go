package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// GET /api/webhook/new_track/<webhook_id>
func apiWebhookNewTrackWebhookIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet { // The method has to be GET

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

	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}
