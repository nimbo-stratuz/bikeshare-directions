package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Routes for resource 'directions'
func Routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", directionsFromTo).Methods("POST")

	return r
}

func directionsFromTo(w http.ResponseWriter, r *http.Request) {

}
