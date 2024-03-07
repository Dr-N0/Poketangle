package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuessHandler(t *testing.T) {
	// Prepare a sample request body
	requestBody := RequestExpect{
		Pokemon:  "Pikachu",
		Question: "What is your favorite move?",
	}
	requestBodyBytes, _ := json.Marshal(requestBody)

	// Create a request using the test server
	request := httptest.NewRequest("POST", "/guess", bytes.NewReader(requestBodyBytes))
	request.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to capture the response
	responseRecorder := httptest.NewRecorder()

	// Call the guessHandler directly (without starting a server) passing the ResponseRecorder and Request
	guessHandler(responseRecorder, request)

	// Check the response status code
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body
	var responseBody map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	expectedResponse := map[string]string{"response": "You asked: What is your favorite move? about Pikachu"}
	if !isEqual(responseBody, expectedResponse) {
		t.Errorf("Handler returned unexpected body: got %v want %v", responseBody, expectedResponse)
	}
}

func isEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if bValue, exists := b[key]; !exists || bValue != value {
			return false
		}
	}
	return true
}
