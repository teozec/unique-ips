package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teozec/unique-ips/internal/services"
)

// Test allowed and not allowed methods for logs handler
func TestLogsAllowedMethods(t *testing.T) {
	uniqueIpCalculator := services.NewUniqueIpCalculator()
	handler := handleLogs(uniqueIpCalculator)

	for _, method := range []string{"GET", "HEAD", "OPTIONS", "TRACE", "PUT", "DELETE", "PATCH", "CONNECT"} {
		req := httptest.NewRequest(method, "/logs", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d; got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	}

	req := httptest.NewRequest("POST", "/logs", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code == http.StatusMethodNotAllowed {
		t.Errorf("unexpected status %d", http.StatusMethodNotAllowed)
	}
}

// Test allowed and not allowed methods for visitors handler
func TestVisitorsAllowedMethods(t *testing.T) {
	uniqueIpCalculator := services.NewUniqueIpCalculator()
	handler := handleVisitors(uniqueIpCalculator)

	for _, method := range []string{"HEAD", "OPTIONS", "TRACE", "POST", "PUT", "DELETE", "PATCH", "CONNECT"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/visitors", nil)
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d; got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	}

	req := httptest.NewRequest("GET", "/logs", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code == http.StatusMethodNotAllowed {
		t.Errorf("unexpected status %d", http.StatusMethodNotAllowed)
	}
}

// Test that visitor handler responds in the correct format ({"count": 5})
func TestVisitorsResponseFormat(t *testing.T) {
	uniqueIpCalculator := services.NewUniqueIpCalculator()
	handler := handleVisitors(uniqueIpCalculator)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/visitors", nil)
	handler.ServeHTTP(rec, req)

	body, _ := io.ReadAll(rec.Body)
	var response map[string]int

	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("error in unmarshalling response to string-int key-value pairs: %s", err)
	}

	if _, ok := response["count"]; !ok {
		t.Errorf("response does not contain \"count\" key")
	}
}
