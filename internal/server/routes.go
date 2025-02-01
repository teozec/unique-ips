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
	mux.Handle("/", http.NotFoundHandler())
	return mux
}

func handleLogs(uniqueIpCalculator services.UniqueIpCalculator) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "POST":
				var err error
				body, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

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
