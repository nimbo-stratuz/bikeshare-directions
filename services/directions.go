package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/service"

	"github.com/nimbo-stratuz/bikeshare-directions/models"
)

type MapQuestService struct {
	url    string
	client *http.Client
}

func NewMapQuestService() MapQuestService {
	apiKey, err := service.Config.Get("maps", "api", "key")
	if err != nil {
		log.Panicln("API key not set")
	}

	url := fmt.Sprintf("https://www.mapquestapi.com/directions/v2/route?key=%s", apiKey)

	return MapQuestService{
		url: url,
		client: &http.Client{
			Timeout: time.Millisecond * 2500,
		},
	}
}

// DirectionsFromTo ...
func (mq *MapQuestService) DirectionsFromTo(from string, to string) models.Directions {

	directionsBody := models.DirectionsRequest{
		Locations: []string{from, to},
		Options: models.DirectionsRequestOptions{
			RouteType: "bicycle",
			Unit:      "k",
		},
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(directionsBody); err != nil {
		log.Panicln("Marshalling failed")
	}

	resp, err := mq.client.Post(mq.url, "application/json; charset=utf-8", body)
	if err != nil {
		log.Panicln("Could not do POST to the maps API")
	}
	defer resp.Body.Close()

	route := models.Directions{}
	if err := json.NewDecoder(resp.Body).Decode(&route); err != nil {

	}

	if route.Info.Statuscode != 0 {
		log.Panicln("API error {}: {}", route.Info.Statuscode, route.Info.Messages)
	}

	return route
}
