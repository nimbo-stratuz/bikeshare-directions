package main

// v0.0.0
import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/config"
	etcd2 "go.etcd.io/etcd/client"

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

	startupConf, err := config.New(
		envConf,
		yamlConf,
	)
	if err != nil {
		log.Fatal(err)
	}

	etcdURL, err := startupConf.Get("config", "etcd", "url")
	if err != nil {
		log.Fatal("config.etcd.url not specified")
	}

	etcd2Conf, err := config.NewEtcd2Config(
		etcd2.Config{
			Endpoints:               []string{etcdURL},
			Transport:               etcd2.DefaultTransport,
			HeaderTimeoutPerRequest: time.Second,
		},
	)
	defer etcd2Conf.Close()
	if err != nil {
		log.Fatal(err)
	}

	multiConf, err := config.New(
		// highest priority
		envConf,
		etcd2Conf,
		yamlConf,
		// lowest priority
	)
	defer multiConf.Close()
	if err != nil {
		log.Fatal(err)
	}

	v4, err := etcd2Conf.Get("foo")
	if err != nil {
		if err == context.Canceled {
			log.Fatal("Etcd2Conf | Canceled: " + context.Canceled.Error())
			// ctx is canceled by another routine
		} else if err == context.DeadlineExceeded {
			log.Fatal("Etcd2Conf | DeadlineExceeded: " + context.DeadlineExceeded.Error())
			// ctx is attached with a deadline and it exceeded
		} else if cerr, ok := err.(*etcd2.ClusterError); ok {
			log.Fatal("Etcd2Conf | ClusterError: " + cerr.Error())
			// process (cerr.Errors)
		} else {
			log.Fatal("Etcd2Conf | OTHER: " + err.Error())
			// bad cluster endpoints, which are not etcd servers
		}
	}
	log.Printf("GET: %v\n", v4)

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

	port, err := startupConf.Get("port")
	if err != nil {
		log.Println("PORT not specified")
	}
	log.Fatal(http.ListenAndServe(":"+port, r))

	time.Sleep(5 * time.Second)
}
