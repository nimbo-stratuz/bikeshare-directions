package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/config"
	etcd2 "go.etcd.io/etcd/client"
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

	yamlConf, err := config.NewYamlFileConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	envConf := config.NewEnvConfig()

	etcdURL, err := envConf.Get("ETCD", "URL")
	if err != nil {
		log.Println("ETCD_URL not specified")
	}

	log.Println(etcdURL)

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

	etcd2Conf, err := config.NewEtcd2Config(
		"",
		// fmt.Sprintf("/%s/%s/", app, instanceID),
		etcd2.Config{
			Endpoints:               []string{etcdURL},
			Transport:               etcd2.DefaultTransport,
			HeaderTimeoutPerRequest: time.Second,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	multiConf, err := config.New(
		// highest priority
		envConf,
		etcd2Conf,
		etcdConf,
		yamlConf,
		// lowest priority
	)
	if err != nil {
		log.Fatal(err)
	}
	defer multiConf.Close()

	v5, err := etcd2Conf.Put("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v5)

	v4, err := etcd2Conf.Get("foo")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v4)

	if _, err := etcdConf.Put("maps/api/key", "1234"); err != nil {
		log.Fatal(err)
	}

	// maps.api.key
	v1, err := multiConf.Get("maps", "api", "key")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v1)

	v2, err := multiConf.GetInt("what", "is", "this")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v2)

	r := Routes()

	port, err := envConf.Get("PORT")
	if err != nil {
		log.Println("PORT not specified")
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
