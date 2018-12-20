package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/config"
	etcd3 "go.etcd.io/etcd/clientv3"

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

		middleware.RequestID,
		middleware.Logger,
	)

	r.Mount("/", api.Routes())

	return r
}

func main() {

	etcdConf, err := config.New(etcd3.Config{
		Endpoints:   []string{getEnv("ETCD_URL", "localhost:2379")},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	if _, err := etcdConf.Put("/some/random/path", "1234"); err != nil {
		log.Fatal(err)
	}

	v1, err := etcdConf.Get("/some/random/path")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v1)

	v2, err := etcdConf.GetInt("/some/random/path")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v2)

	r := Routes()
	log.Fatal(http.ListenAndServe(":"+getEnv("PORT", "8080"), r))
}

func getEnv(k string, def string) string {
	v := os.Getenv(k)

	if v != "" {
		return v
	}

	return def
}
