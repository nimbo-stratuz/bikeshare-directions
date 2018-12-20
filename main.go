package main

import (
	"fmt"
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

	"github.com/google/uuid"
)

var (
	app        = "bikeshare-directions"
	instanceID = uuid.New().String()
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

	etcdConf, err := config.NewEtcdConfig(
		fmt.Sprintf("/%s/%s/", app, instanceID),
		etcd3.Config{
			Endpoints:   []string{getEnv("ETCD_URL", "localhost:2379")},
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	envConf := config.NewEnvConfig()

	multiConf := config.New(
		etcdConf,
		envConf,
	)
	defer multiConf.Close()

	if _, err := etcdConf.Put("MAPS_API_KEY", "1234"); err != nil {
		log.Fatal(err)
	}

	v1, err := multiConf.Get("MAPS_API_KEY")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v1)

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
