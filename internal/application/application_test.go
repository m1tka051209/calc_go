package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalcHandler_ValidExpression(t *testing.T) {
	// Создаём тестовый запрос с допустимым выражением.
	reqBody := `{"expression": "2+2"}`
	req, err := http.NewRequest("POST", "/", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаём ResponseRecorder для записи ответа.
	rr := httptest.NewRecorder()

	// Вызоваем CalcHandler с запросом и ResponseRecorder.
	CalcHandler(rr, req)

	// Проверим код состояния.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверим тело ответа.
	expected := "result: 4.000000"
	if got := rr.Body.String(); !strings.Contains(got, expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
	}
}

func TestCalcHandler_InvalidExpression(t *testing.T) {
	// Создаём тестовый запрос с недопустимым выражением.
	reqBody := `{"expression": "2+a"}`
	req, err := http.NewRequest("POST", "/", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаём ResponseRecorder для записи ответа.
	rr := httptest.NewRecorder()

	// Вызываем CalcHandler с запросом и ResponseRecorder.
	CalcHandler(rr, req)

	// Проверяем код состояния.
	if status := rr.Code; status != http.StatusOK { //Status will be 200 since we handle the error case.
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем тело ответа.
	expected := "err: invalid expression"
    if got := rr.Body.String(); !strings.Contains(got, expected) { // Use strings.Contains instead of direct equality
        t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
    }
}

func TestRunServer_Basic(t *testing.T) {
  mux := http.NewServeMux()
	mux.HandleFunc("/", CalcHandler)

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	// Выполняем тестовый запрос.
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

	// Читаем ответ.
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	got := buf.String()

	// Проверяем ответ
	expected := "result: 15.000000"
	if !strings.Contains(got, expected) {
		t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected a 200 status code, but got %v", resp.StatusCode)
	}
}
