package main

import (
	"fmt"
	"log"
	"net/http"
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

	envConf := config.NewEnvConfig()

	etcdURL, err := envConf.Get("ETCD_URL")
	if err != nil {
		log.Println("ETCD_URL not specified")
	}

	etcdConf, err := config.NewEtcdConfig(
		fmt.Sprintf("/%s/%s/", app, instanceID),
		etcd3.Config{
			Endpoints:   []string{etcdURL},
			DialTimeout: 5 * time.Second,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	multiConf := config.New(
		envConf,  // highest priority
		etcdConf, // lowest priority
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

	port, err := envConf.Get("PORT")
	if err != nil {
		log.Println("ETCD_URL not specified")
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
