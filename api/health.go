package api

import (
	"net/http"
	"os"

	"github.com/go-chi/render"
)

// HealthCheckResponse is a microprofile-like /health response
type HealthCheckResponse struct {
	Outcome string           `json:"outcome"`
	Checks  []SubHealthCheck `json:"checks"`
}

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
	url = "https://www.mapquestapi.com/directions/v2/route?key=" + os.Getenv("MAPS_API_KEY")
)

const (
	stateUp   = "UP"
	stateDown = "DOWN"
)

// HealthCheck is a basic Healthcheck
func HealthCheck(w http.ResponseWriter, r *http.Request) {
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
