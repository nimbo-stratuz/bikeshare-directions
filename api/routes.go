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

	route := services.DirectionsFromTo(fromTo.From, fromTo.To)

	render.Render(w, r, &route)
}

// ErrInvalidRequest creates an ErrResponse for 400 Bad Request
func ErrInvalidRequest(err error) render.Renderer {
	return &models.ErrResponse{
		StatusCode: 400,
		ErrorText:  "Invalid request",
	}
}
