package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nimbo-stratuz/bikeshare-directions/api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Routes for /v1
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		render.SetContentType(render.ContentTypeJSON),

		middleware.Logger,
		middleware.RequestID,
	)

	r.Mount("/", api.Routes())

	return r
}

func main() {
	r := Routes()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
