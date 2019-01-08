package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/middleware"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/render"

	"github.com/nimbo-stratuz/bikeshare-directions/service"

	"github.com/nimbo-stratuz/bikeshare-directions/models"
)

// DirectionsFromTo ...
func DirectionsFromTo() http.HandlerFunc {

	apiKey, err := service.Config.Get("maps", "api", "key")
	if err != nil {
		log.Fatal(err)
	}

	mapsURL := fmt.Sprintf("https://www.mapquestapi.com/directions/v2/route?key=%s", apiKey)

	client := &http.Client{
		Timeout: time.Millisecond * 2500,
	}

	return func(w http.ResponseWriter, r *http.Request) {

		var (
			bycicle models.Bicycle    = models.Bicycle{}
			route   models.Directions = models.Directions{}
		)

		fromTo := &models.FromTo{}

		if err := render.Bind(r, fromTo); err != nil {
			render.Render(w, r, ErrBadRequest("From/to not specified"))
			return
		}

		// GET route
		{
			directionsBody := models.DirectionsRequest{
				Locations: []string{fromTo.From, fromTo.To},
				Options: models.DirectionsRequestOptions{
					RouteType: "bicycle",
					Unit:      "k",
				},
			}

			body := new(bytes.Buffer)
			if err := json.NewEncoder(body).Encode(directionsBody); err != nil {
				log.Print("MapsQuest request/Encode:", err)
				render.Render(w, r, ErrServerError())
				return
			}

			resp, err := client.Post(mapsURL, "application/json; charset=utf-8", body)
			if err != nil {
				log.Print(err)
				render.Render(w, r, Err(503, "Maps API unavaliable"))
				return
			}
			defer resp.Body.Close()

			if err := json.NewDecoder(resp.Body).Decode(&route); err != nil {
				log.Println("MapsQuest response/Decode:", err)
				render.Render(w, r, ErrServerError())
				return
			}

			if route.Info.Statuscode != 0 {
				log.Printf("MapsQuest API Error code = %d\n", route.Info.Statuscode)
				render.Render(w, r, Err(503, "Maps API Error"))
				return
			}
		}

		// GET closest bicycle
		{
			catalogueURLString, err := service.Discovery.Discover("bikeshare-catalogue", service.GetEnv(), "1.0.0")
			if err != nil {
				log.Println(err)
				render.Render(w, r, ErrServerError())
				return
			}

			catalogueURL, err := url.Parse(catalogueURLString)
			if err != nil {
				log.Println(err)
				render.Render(w, r, ErrServerError())
				return
			}

			query := catalogueURL.Query()

			query.Set("latitude", fmt.Sprint(route.Route.Locations[0].LatLng.Lat))
			query.Set("longitude", fmt.Sprint(route.Route.Locations[0].LatLng.Lng))

			catalogueURL.RawQuery = query.Encode()

			req, err := http.NewRequest("GET", catalogueURL.String(), nil)
			if err != nil {
				log.Panic(err)
			}

			req.Header.Set("X-Request-ID", fmt.Sprint(r.Context().Value(middleware.RequestIDKey)))

			resp, err := client.Do(req)
			if err != nil {
				log.Panic(err)
			}
			defer resp.Body.Close()

			json.NewDecoder(resp.Body).Decode(&bycicle)
		}

		render.Render(w, r, &models.DirectionsWithBicycle{
			Bicycle:    &bycicle,
			Directions: &route,
		})
	}
}
