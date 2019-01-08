package api

import (
	"github.com/go-chi/chi"
	"github.com/nimbo-stratuz/bikeshare-directions/handlers"
)

// Routes for resource 'directions'
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Route("/directions", func(r chi.Router) {
			r.Post("/", handlers.DirectionsFromTo())
		})
	})

	r.Route("/health", func(r chi.Router) {
		r.Get("/", HealthCheck)
		r.Head("/", HealthCheck)
	})

	return r
}
