package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/teozec/unique-ips/internal/services"
)

func (s *Server) RegisterRoutes(uniqueIpCalculator services.UniqueIpCalculator) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/logs", handleLogs(uniqueIpCalculator))
	mux.Handle("/visitors", handleVisitors(uniqueIpCalculator))
	mux.Handle("/", http.NotFoundHandler())
	return mux
}

// Receives POST requests in the following format:
// { "timestamp": "2020-06-24T15:27:00.123456Z", "ip": "83.150.59.250", "url": ... }
// and sends the IP addresses to the uniqueIpCalculator logs.
func handleLogs(uniqueIpCalculator services.UniqueIpCalculator) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				// Read the request body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				// Try to unmarshal the json
				var log map[string]string
				if err = json.Unmarshal([]byte(body), &log); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				uniqueIpCalculator.LogIp(log["ip"])

			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)
}

// Receives GET requests and responds with the number of unique IP addresses that have been logged.
// Responses are in the format { "count": 5 }
func handleVisitors(uniqueIpCalculator services.UniqueIpCalculator) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET", "":
				// Build the response
				count := uniqueIpCalculator.GetUniqueIpNumber()
				data := map[string]int{"count": count}

				// Marshal the response to json
				response, err := json.Marshal(data)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				w.Write(response)

			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)
}
