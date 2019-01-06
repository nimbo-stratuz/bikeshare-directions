package api

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/nimbo-stratuz/bikeshare-directions/models"

	"github.com/go-chi/chi"
	"github.com/nimbo-stratuz/bikeshare-directions/services"
)

// Routes for resource 'directions'
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Route("/directions", func(r chi.Router) {
			r.Post("/", findDirections)
		})
	})

	r.Route("/health", func(r chi.Router) {
		r.Get("/", HealthCheck)
		r.Head("/", HealthCheck)
	})

	return r
}

func findDirections(w http.ResponseWriter, r *http.Request) {

	fromTo := &models.FromTo{}

	if err := render.Bind(r, fromTo); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	directionsAPI := services.NewMapQuestService()
	catalogueAPI := services.NewBikeshareCatalogueService()

	route := directionsAPI.DirectionsFromTo(fromTo.From, fromTo.To)

	startLatitude := route.Route.Locations[0].LatLng.Lat
	startLongitude := route.Route.Locations[0].LatLng.Lng

	bicycle := catalogueAPI.ClosestBicycle(startLatitude, startLongitude, r)

	result := models.DirectionsWithBicycle{
		Bicycle:    &bicycle,
		Directions: &route,
	}

	render.Render(w, r, &result)
}

// ErrInvalidRequest creates an ErrResponse for 400 Bad Request
func ErrInvalidRequest(err error) render.Renderer {
	return &models.ErrResponse{
		StatusCode: 400,
		ErrorText:  "Invalid request",
	}
}
