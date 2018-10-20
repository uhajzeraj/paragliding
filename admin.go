package main

import (
	"fmt"
	"net/http"
)

func adminAPITrackCountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet { // GET request

		// Connect to DB
		conn := mongoConnect()

		// Get the tracks
		trackCount := len(getAllTracks(conn, false))

		fmt.Fprintln(w, trackCount)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}

func adminAPITracksDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete { // DELETE request

		// Connect to DB
		conn := mongoConnect()

		// Get the tracks
		trackCount := len(getAllTracks(conn, false))

		// Delete all the tracks
		deleteAllTracks(conn)

		fmt.Fprintln(w, trackCount)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404 Not Found
	}
}
