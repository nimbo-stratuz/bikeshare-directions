package main

import (
	"encoding/json"
	"net/http"
)

// HealthCheckResponse is a microprofile-like /health response
type HealthCheckResponse struct {
	Outcome string           `json:"outcome"`
	Checks  []SubHealthCheck `json:"checks"`
}

// SubHealthCheck represents underlying health checks
type SubHealthCheck struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

// HealthCheck is a basic Healthcheck
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	hc := HealthCheckResponse{
		Outcome: "OK",
		Checks: []SubHealthCheck{
			SubHealthCheck{
				Name:  "GoogleMapsHealthCheck",
				State: "UP",
			},
		},
	}

	json, err := json.Marshal(hc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
