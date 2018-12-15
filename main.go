package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nimbo-stratuz/bikeshare-directions/api"

	"github.com/gorilla/mux"
)

// Routes for /v1
func Routes() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/directions", api.Routes())

	return r
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/health", HealthCheck).Methods("GET")

	r.Handle("/v1", Routes())

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "200 OK")
	}).Methods("GET")

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
