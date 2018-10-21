package main

import (
	"encoding/json"
	"net/http"
)

// // Webhook - Handle the decoded JSON body from the POST request
// type Webhook struct {
// 	WebhookURL      string `json:"webhook_url"`
// 	MinTriggerValue int    `json:"min_trigger_value"`
// 	WebhookID       string `json:"webhook_id"`
// }

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

		insertUpdateWebhook(data)

	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}
