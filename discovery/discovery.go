package discovery

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nimbo-stratuz/bikeshare-directions/config"
	etcd2 "go.etcd.io/etcd/client"
)

type ServiceDiscovery interface {
	Register()
	Discover(service, env, version string) string
	Close()
}

type discovery struct {
	instanceID string

	config config.Config
	kapi   etcd2.KeysAPI

	refresherChan chan bool
}

var (
	ttl     = 10 * time.Second
	refresh = ttl / 2
)

func New(instanceID string, cfg config.Config) (ServiceDiscovery, error) {

	etcdURL, err := cfg.Get("discovery", "etcd", "url")
	if err != nil {
		log.Fatal(err)
	}

	etcd2Client, err := etcd2.New(
		etcd2.Config{
			Endpoints:               []string{etcdURL},
			Transport:               etcd2.DefaultTransport,
			HeaderTimeoutPerRequest: time.Second,
		},
	)

	if err != nil {
		return nil, err
	}

	return &discovery{
		instanceID:    instanceID,
		config:        cfg,
		kapi:          etcd2.NewKeysAPI(etcd2Client),
		refresherChan: make(chan bool),
	}, nil
}

func (d *discovery) Close() {
	log.Println("Closing ServiceDiscovery")
	d.Deregister()
	d.refresherChan <- true
}

func (d *discovery) Register() {

	url, err := d.config.Get("server", "baseurl")
	if err != nil {
		log.Fatal(err)
	}

	// Create directory with TTL
	{
		path := d.genPathInstance()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		_, err := d.kapi.Set(ctx, path, "", &etcd2.SetOptions{
			TTL: ttl,
			Dir: true,
		})
		defer cancel()
		if err != nil {
			log.Fatal("Could not register service (create dir)", err)
			return
		}
	}

	// Set value
	{
		path := d.genPathInstanceURL()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		_, err := d.kapi.Set(ctx, path, url, nil)
		defer cancel()
		if err != nil {
			log.Fatal("Could not register service (set value)", err)
			return
		}
	}

	// Keep refreshing
	go func() {
		log.Printf("Refreshing every %s\n", refresh)
		for {
			time.Sleep(refresh)

			log.Println("Refreshing service in discovery")

			path := d.genPathInstance()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			_, err := d.kapi.Set(ctx, path, "", &etcd2.SetOptions{
				TTL:       ttl,
				Dir:       true,
				PrevExist: etcd2.PrevExist,
			})
			defer cancel()
			if err != nil {
				log.Fatal("Could not refresh service in service discovery", err)
			}

			select {
			case <-d.refresherChan:
				log.Println("No longer refreshing")
				break
			default:
				continue
			}
		}
	}()
}

func (d *discovery) Deregister() {

	path := d.genPathInstance()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err := d.kapi.Delete(ctx, path, &etcd2.DeleteOptions{
		Dir:       true,
		Recursive: true,
	})
	defer cancel()
	if err != nil {
		log.Fatal("Could not deregister service", err)
		return
	}
}

func (d *discovery) Discover(name, env, version string) string {

	instances := d.list(d.genPath())
	if len(instances) <= 0 {
		log.Println("Cannot discover service " + name)
		return ""
	}

	idx := rand.Intn(len(instances))

	path := instances[idx] + "/url"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := d.kapi.Get(ctx, path, nil)
	defer cancel()
	if err != nil {
		log.Printf("Cannot discover service %s | %s | %s: %s\n", name, env, version, err.Error())
		return ""
	}

	return resp.Node.Value
}

func (d *discovery) list(dir string) []string {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := d.kapi.Get(ctx, dir, nil)
	defer cancel()
	if err != nil {
		log.Println(err)
		return []string{}
	}

	keys := []string{}

	log.Printf("%+v\n", resp.Node.Nodes)

	for _, node := range resp.Node.Nodes {
		keys = append(keys, node.Key)
	}

	return keys
}

func (d *discovery) genPath() string {
	env, err := d.config.Get("env")
	if err != nil {
		log.Fatal(err)
	}

	name, err := d.config.Get("name")
	if err != nil {
		log.Fatal(err)
	}

	version, err := d.config.Get("version")
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("/environments/%s/services/%s/%s/instances", env, name, version)
}

func (d *discovery) genPathInstance() string {
	return fmt.Sprintf("%s/%s", d.genPath(), d.instanceID)
}

func (d *discovery) genPathInstanceURL() string {
	return fmt.Sprintf("%s/url", d.genPathInstance())
}
