package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nimbo-stratuz/bikeshare-directions/models"
)

var (
	url = "http://www.mapquestapi.com/directions/v2/route?key=" + os.Getenv("MAPS_API_KEY")
)

// DirectionsFromTo ...
func DirectionsFromTo(from string, to string) models.Directions {
	directionsBody := models.DirectionsRequest{
		Locations: []string{from, to},
		Options: models.DirectionsRequestOptions{
			RouteType: "bicycle",
			Unit:      "k",
		},
	}

	body, err := json.Marshal(directionsBody)
	if err != nil {
		log.Panicln("Marshalling failed")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Panicln("Couldn't create request")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln("Could do POST to the maps API")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("could not read resp.Body")
	}

	fmt.Println(string(respBody))

	route := models.Directions{}

	err = json.Unmarshal(respBody, &route)
	if err != nil {
		log.Panicln("could not unmarshall json response")
	}

	if route.Info.Statuscode != 0 {
		log.Panicln("API error {}: {}", route.Info.Statuscode, route.Info.Messages)
	}

	return route
}
