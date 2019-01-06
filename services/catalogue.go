package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/service"

	"github.com/nimbo-stratuz/bikeshare-directions/models"

	"github.com/go-chi/chi/middleware"
)

type BikeshareCatalogueService struct {
	baseURL string
	client  *http.Client
}

func NewBikeshareCatalogueService() BikeshareCatalogueService {

	baseURL, err := service.Discovery.Discover("bikeshare-catalogue", service.GetEnv(), "1.0.0")
	if err != nil {
		log.Panicln("Cannot discover bikeshare-catalogue")
	}

	return BikeshareCatalogueService{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Millisecond * 2500,
		},
	}
}

// DirectionsFromTo ...
func (bc *BikeshareCatalogueService) ClosestBicycle(latitude, longitude float64, r *http.Request) models.Bicycle {

	apiURL, err := url.Parse(bc.baseURL + "/v1/bicycles")
	if err != nil {
		log.Panic(err)
	}

	query := apiURL.Query()

	query.Set("latitude", fmt.Sprint(latitude))
	query.Set("longitude", fmt.Sprint(longitude))

	apiURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", apiURL.String(), nil)
	if err != nil {
		log.Panic(err)
	}

	req.Header.Set("X-Request-ID", fmt.Sprint(r.Context().Value(middleware.RequestIDKey)))

	resp, err := bc.client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	var bicycle models.Bicycle
	json.NewDecoder(resp.Body).Decode(&bicycle)

	return bicycle
}
