package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_apiHandler(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/paragliding/api", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(apiHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}

func Test_redirectHandler(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("GET", "/paragliding", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusSeeOther { // It should be 303 (SeeOther)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}
