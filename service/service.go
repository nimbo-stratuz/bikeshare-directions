package service

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nimbo-stratuz/bikeshare-directions/config"
	"github.com/nimbo-stratuz/bikeshare-directions/discovery"

	etcd2 "go.etcd.io/etcd/client"
)

var (
	// InstanceID ...
	InstanceID = uuid.New().String()

	// Config ...
	Config config.Config

	// Discovery ...
	Discovery discovery.ServiceDiscovery
)

func init() {
	initConfig()
	initDiscovery()
}

func initConfig() {
	log.Println("Initializing Config")

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
	if err != nil {
		log.Fatal(err)
	}

	Config, err = config.New(
		// highest priority
		envConf,
		etcd2Conf,
		yamlConf,
		// lowest priority
	)
	if err != nil {
		log.Fatal(err)
	}
}

func initDiscovery() {
	log.Println("Initializing Discovery")

	var err error
	Discovery, err = discovery.New(InstanceID, Config)
	if err != nil {
		log.Fatal(err)
	}

	Discovery.Register()

	log.Println(Discovery.Discover("bikeshare-directions", "dev", "1.0.0"))
}
