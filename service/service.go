package service

import (
	"time"

	log "github.com/sirupsen/logrus"

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
	initLogging()
	initConfig()
	initDiscovery()
}

func Close() {
	Discovery.Close()
	Config.Close()
}

func initLogging() {
	log.SetLevel(log.DebugLevel)
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

	// Initialize Discovery
	{
		var err error
		Discovery, err = discovery.New(InstanceID, Config)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Register the service
	{
		err := Discovery.Register()

		for i := 0; i < 2 && err != nil; i++ {
			log.Debugf("Retrying to register in 2 seconds... (Retry #%d)", i+1)
			<-time.After(2 * time.Second)
			err = Discovery.Register()
		}

		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetName() string {
	env, err := Config.Get("name")
	if err != nil {
		log.Panic("Name not set")
	}

	return env
}

func GetEnv() string {
	env, err := Config.Get("env")
	if err != nil {
		log.Panic("Env not set")
	}

	return env
}
