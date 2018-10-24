package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_apiIgcHandler(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/paragliding/api/igc", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(apiIgcHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}

func Test_apiIgcIDHandler(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/paragliding/api/igc/igc1", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(apiIgcIDHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code | 404 because there are no records added
	if resRecorder.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusNotFound)
	}
}

func Test_apiIgcIDFieldHandler(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/paragliding/api/igc", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(apiIgcIDFieldHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code | 404 because there are no records added
	if resRecorder.Code != http.StatusNotFound { // It should be 404 (Not Found)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}
