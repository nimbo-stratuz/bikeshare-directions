package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/render"
	"github.com/nimbo-stratuz/bikeshare-directions/service"
)

// HealthCheckResponse is a microprofile-like /health response
type HealthCheckResponse struct {
	Outcome string           `json:"outcome"`
	Checks  []SubHealthCheck `json:"checks"`
}

// Render sets HTTP Status to 503 if outcome equals DOWN
func (hcr *HealthCheckResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if hcr.Outcome == stateDown {
		render.Status(r, http.StatusServiceUnavailable)
	}

	return nil
}

// SubHealthCheck represents underlying health checks
type SubHealthCheck struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

var (
	url = ""
)

const (
	stateUp   = "UP"
	stateDown = "DOWN"
)

// HealthCheck is a basic Healthcheck
func HealthCheck(w http.ResponseWriter, r *http.Request) {

	if url == "" {
		apiKey, err := service.Config.Get("maps", "api", "key")
		if err != nil {
			log.Panicln("API key not set")
		}

		url = fmt.Sprintf("https://www.mapquestapi.com/directions/v2/route?key=%s", apiKey)
	}

	checks := []SubHealthCheck{
		doMapQuestHealthCheck(),
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

	render.Render(w, r, &hc)
}

func mapQuestHealthCheck(isUp bool) SubHealthCheck {

	var state string
	if isUp {
		state = stateUp
	} else {
		state = stateDown
	}

	return SubHealthCheck{
		Name:  "MapQuestHealthCheck",
		State: state,
	}
}

func doMapQuestHealthCheck() SubHealthCheck {

	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return mapQuestHealthCheck(false)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return mapQuestHealthCheck(false)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return mapQuestHealthCheck(false)
	}

	return mapQuestHealthCheck(true)
}
