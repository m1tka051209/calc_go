package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalcHandler_ValidExpression(t *testing.T) {
	// Create a test request with a valid expression.
	reqBody := `{"expression": "2+2"}`
	req, err := http.NewRequest("POST", "/", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// Call the CalcHandler with the request and ResponseRecorder.
	CalcHandler(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body.
	expected := "result: 4.000000"
	if got := rr.Body.String(); !strings.Contains(got, expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestCalcHandler_InvalidExpression(t *testing.T) {
	// Create a test request with an invalid expression.
	reqBody := `{"expression": "2+a"}`
	req, err := http.NewRequest("POST", "/", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to capture the response.
	rr := httptest.NewRecorder()

	// Call the CalcHandler with the request and ResponseRecorder.
	CalcHandler(rr, req)

	// Check the status code.
	if status := rr.Code; status != http.StatusOK { //Status will be 200 since we handle the error case.
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body.
	expected := "err: invalid expression"
    if got := rr.Body.String(); !strings.Contains(got, expected) { // Use strings.Contains instead of direct equality
        t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
    }
}

func TestRunServer_Basic(t *testing.T) {
    // Create a new http serve mux, and add the CalcHandler to it at the desired path.
	mux := http.NewServeMux()
	mux.HandleFunc("/", CalcHandler)

    // Create the test server using the created serve mux, so we can route requests to our correct path.
	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	// Perform a test request.
	reqBody := `{"expression": "5*3"}`
	req, err := http.NewRequest("POST", testServer.URL, strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	//Read the response.
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	got := buf.String()

	// Check the response
	expected := "result: 15.000000"
	if !strings.Contains(got, expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 status code, but got %v", resp.StatusCode)
	}
}
