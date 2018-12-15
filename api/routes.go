package api

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/nimbo-stratuz/bikeshare-directions/models"

	"github.com/go-chi/chi"
	"github.com/nimbo-stratuz/bikeshare-directions/services"
)

// Routes for resource 'directions'
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {

		r.Route("/directions", func(r chi.Router) {

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {

				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Panicln("Canot read request body")
				}

				fromTo := models.FromTo{}

				unmarshallErr := json.Unmarshal(body, &fromTo)
				if unmarshallErr != nil {
					log.Panicln("Cannot unmarshal json")
				}

				fmt.Printf("%v\n", fromTo)

				route := services.DirectionsFromTo(fromTo.From, fromTo.To)

				json.NewEncoder(w).Encode(route)
			})

		})
	})

	r.Get("/health", HealthCheck)

	return r
}
