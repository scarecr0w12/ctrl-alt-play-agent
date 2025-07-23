package health

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthStatus represents the health status of the agent
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	NodeID    string    `json:"nodeId"`
	Uptime    string    `json:"uptime"`
	Connected bool      `json:"connected"`
}

// Server provides health check endpoints
type Server struct {
	startTime time.Time
	nodeID    string
	version   string
	connected bool
}

// NewServer creates a new health check server
func NewServer(nodeID, version string) *Server {
	return &Server{
		startTime: time.Now(),
		nodeID:    nodeID,
		version:   version,
		connected: false,
	}
}

// SetConnectionStatus updates the connection status
func (s *Server) SetConnectionStatus(connected bool) {
	s.connected = connected
}

// Handler returns the HTTP handler for health checks
func (s *Server) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
			Version:   s.version,
			NodeID:    s.nodeID,
			Uptime:    time.Since(s.startTime).String(),
			Connected: s.connected,
		}

		if !s.connected {
			status.Status = "degraded"
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(status); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}

// StartServer starts the health check server on the specified port
func (s *Server) StartServer(port string) error {
	http.HandleFunc("/health", s.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/health", http.StatusMovedPermanently)
	})

	return http.ListenAndServe(":"+port, nil)
}
