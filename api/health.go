package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
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

var (
	url = "https://www.mapquestapi.com/directions/v2/route?key=" + os.Getenv("MAPS_API_KEY")
)

const (
	stateUp   = "UP"
	stateDown = "DOWN"
)

// HealthCheck is a basic Healthcheck
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	checks := []SubHealthCheck{
		mapQuestHealthCheck(),
	}

	state := stateUp

	for _, chk := range checks {
		if chk.State == stateDown {
			state = stateDown
			break
		}
	}

	hc := HealthCheckResponse{
		Outcome: state,
		Checks:  checks,
	}

	json, err := json.Marshal(hc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if state == stateDown {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func mapQuestHealthCheck() SubHealthCheck {

	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		log.Panicln("Couldn't create request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("Couldn't do OPTIONS to the maps API", err.Error())
	}
	defer resp.Body.Close()

	var state string
	log.Println(resp.StatusCode)
	if resp.StatusCode == 200 {
		state = stateUp
	} else {
		state = stateDown
	}

	return (SubHealthCheck{
		Name:  "MapQuestHealthCheck",
		State: state,
	})
}
