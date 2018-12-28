package main

// v0.0.0
import (
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/api"
	"github.com/nimbo-stratuz/bikeshare-directions/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// Routes for /v1
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		render.SetContentType(render.ContentTypeJSON),

		middleware.RequestID,
		middleware.Logger,
	)

	r.Mount("/", api.Routes())

	return r
}

func main() {
	defer service.Config.Close()
	defer service.Discovery.Close()

	// r := Routes()

	// port, err := service.Config.Get("server", "port")
	// if err != nil {
	// 	log.Println("PORT not specified")
	// }
	// log.Fatal(http.ListenAndServe(":"+port, r))

	time.Sleep(10 * time.Second)
}
