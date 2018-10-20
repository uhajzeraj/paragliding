package main

import (
	"encoding/json"
	"net/http"
)

func apiWebhookNewTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost { // The method has to be POST
		var data map[string]interface{} // Save POST body (which is JSON) in the map

		err := json.NewDecoder(r.Body).Decode(&data) // Json decode
		if err != nil {
			return
		}

		// Check if minTriggerValue exists in the map
		if _, ok := data["minTriggerValue"]; ok {

		}

		// Check if webhookURL exists in the map
		if _, ok := data["webhookURL"]; ok {

		} else {
			return
		}

	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}
