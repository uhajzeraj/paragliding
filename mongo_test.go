package main

import (
	"testing"
)

func Test_mongoConnect(t *testing.T) {
	if conn := mongoConnect(); conn == nil {
		t.Error("No connection")
	}
}

func Test_urlInMong(t *testing.T) {
	urlExists := urlInMongo(`Some random URL`, mongoConnect().Database("paragliding").Collection("track"))
	if urlExists {
		t.Error("Track should not exist")
	}
}

func Test_trackNameFromURL(t *testing.T) {
	conn := mongoConnect()
	result := trackNameFromURL("url", conn.Database("paragliding").Collection("track"))

	if result != "" {
		t.Error("Name should not exist")
	}
}

func Test_getCounter(t *testing.T) {
	conn := mongoConnect()
	count := getCounter(conn.Database("paragliding"), "trackCounter")

	if count < 1 {
		t.Error("It should be bigger")
	}
}

func Test_getAllTrack(t *testing.T) {
	allTracks := getAllTracks(mongoConnect())

	if len(allTracks) < 0 {
		t.Error("It should be bigger")
	}
}

func Test_getAllWebhooks(t *testing.T) {
	allWebhooks := getAllWebhooks(mongoConnect())

	if len(allWebhooks) < 0 {
		t.Error("It should be bigger")
	}
}

func Test_getTrack(t *testing.T) {
	track := getTrack(mongoConnect(), `url`)

	if track.TrackName != "" {
		t.Error("It should be empty")
	}
}

func Test_getWebhook(t *testing.T) {
	webhook := getWebhook(mongoConnect(), `webhook`)

	if webhook.WebhookID != "" {
		t.Error("It should be empty")
	}
}

func Test_deleteWebhook(t *testing.T) {
	deleteWebhook(mongoConnect(), `noWebhook`)
}
