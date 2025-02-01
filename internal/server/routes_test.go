package server

import (
	"bytes"
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

// End-to-end test to see if the endpoints are wirking as expected
func TestHandlers(t *testing.T) {
	uniqueIpCalculator := services.NewUniqueIpCalculator()
	logsHandler := handleLogs(uniqueIpCalculator)
	visitorsHandler := handleVisitors(uniqueIpCalculator)

	// Test the insertion of 30 distinct IPs.
	ips1 := []string{
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
		"192.168.1.4",
		"192.168.1.5",
		"192.168.1.6",
		"192.168.1.7",
		"192.168.1.8",
		"192.168.1.9",
		"192.168.1.10",
		"192.168.1.11",
		"192.168.1.12",
		"192.168.1.13",
		"192.168.1.14",
		"192.168.1.15",
		"192.168.1.16",
		"192.168.1.17",
		"192.168.1.18",
		"192.168.1.19",
		"192.168.1.20",
		"192.168.1.21",
		"192.168.1.22",
		"192.168.1.23",
		"192.168.1.24",
		"192.168.1.25",
		"192.168.1.26",
		"192.168.1.27",
		"192.168.1.28",
		"192.168.1.29",
		"192.168.1.30",
	}

	var count int
	count = logIPsAndGetUniqueCount(ips1, visitorsHandler, logsHandler)
	if count != 30 {
		t.Errorf("expected 30 unique ips, got %d", count)
	}

	// Insert the same 30 IPs and ensure that the count does not change
	count = logIPsAndGetUniqueCount(ips1, visitorsHandler, logsHandler)
	if count != 30 {
		t.Errorf("expected 30 unique ips, got %d", count)
	}

	// Insert 30 new IPs, of which 10 have already been logged. The unique IPs should now be 50.
	ips2 := []string{
		"192.168.1.31",
		"192.168.1.32",
		"192.168.1.33",
		"192.168.1.34",
		"192.168.1.35",
		"192.168.1.36",
		"192.168.1.37",
		"192.168.1.38",
		"192.168.1.39",
		"192.168.1.40",
		"192.168.1.11",
		"192.168.1.12",
		"192.168.1.13",
		"192.168.1.14",
		"192.168.1.15",
		"192.168.1.16",
		"192.168.1.17",
		"192.168.1.18",
		"192.168.1.19",
		"192.168.1.20",
		"192.168.1.51",
		"192.168.1.52",
		"192.168.1.53",
		"192.168.1.54",
		"192.168.1.55",
		"192.168.1.56",
		"192.168.1.57",
		"192.168.1.58",
		"192.168.1.59",
		"192.168.1.60",
	}

	count = logIPsAndGetUniqueCount(ips2, visitorsHandler, logsHandler)
	if count != 50 {
		t.Errorf("expected 50 unique ips, got %d", count)
	}
}

// Utility function to post some ips and get the count from /visitors
func logIPsAndGetUniqueCount(ips []string, visitorsHandler http.Handler, logsHandler http.Handler) int {
	for _, ip := range ips {
		body, _ := json.Marshal(map[string]string{
			"timestamp": "2020-06-24T15:27:00.123456Z",
			"ip":        ip,
			"url":       "/my/url",
		})
		req := httptest.NewRequest("POST", "/visitors", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		logsHandler.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest("GET", "/visitors", nil)
	rec := httptest.NewRecorder()
	visitorsHandler.ServeHTTP(rec, req)

	var response map[string]int
	body, _ := io.ReadAll(rec.Body)
	json.Unmarshal([]byte(body), &response)

	return response["count"]
}
