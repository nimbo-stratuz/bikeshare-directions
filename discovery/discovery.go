package discovery

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/nimbo-stratuz/bikeshare-directions/config"
	etcd2 "go.etcd.io/etcd/client"
)

// ServiceDiscovery is an interface for registering service with etcd
// and discovering other services
type ServiceDiscovery interface {
	Register() error                                       // Register running service with etcd
	Discover(service, env, version string) (string, error) // Discover url of some env/service/version
	Close()                                                // Close stops refreshing TTL and deregisters the service
}

type discovery struct {
	instanceID string

	config config.Config
	kapi   etcd2.KeysAPI

	refresherChan chan bool // Write to this channel to stop refreshing
}

var (
	ttl     = 10 * time.Second
	refresh = ttl / 2
)

// New creates a new ServiceDiscovery instance
func New(instanceID string, cfg config.Config) (ServiceDiscovery, error) {

	etcdURL, err := cfg.Get("discovery", "etcd", "url")
	if err != nil {
		log.Panic("Service discovery misconfigured:", err)
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
	log.Info("Closing ServiceDiscovery")
	d.refresherChan <- true
	d.Deregister()
}

func (d *discovery) Register() error {

	log.Infof("Registering service with etcd. TTL=%s, refresh=%s", ttl, refresh)

	url, err := d.config.Get("server", "baseurl")
	if err != nil {
		log.Panic("Service discovery: baseurl not set. ", err)
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
			return NewRegisterError(err.Error())
		}
	}

	// Set value
	{
		path := d.genPathInstanceURL()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		_, err := d.kapi.Set(ctx, path, url, nil)
		defer cancel()
		if err != nil {
			return NewRegisterError(err.Error())
		}
	}

	// Keep refreshing
	go func() {
		refresherUUID := uuid.New().String()

		log.Infof("Refreshing every %s [Refresher: %s]", refresh, refresherUUID)
		defer log.Infof("No longer refreshing service [Refresher: %s]", refresherUUID)

		for {
			select {

			case <-d.refresherChan:
				return

			case <-time.After(refresh):
				err := d.refresh()

				// Retry loop
				if err != nil {
					for i := 0; i < 5 && err != nil; i++ {
						if etcd2.IsKeyNotFound(err) {
							d.Register()
							return
						}

						log.Warn(err.Error())
						log.Warnf("Refreshing failed. Retrying in 2 seconds (Retry #%d)", i+1)
						<-time.After(2 * time.Second)
						err = d.refresh()
					}

					if err != nil {
						log.Fatal("Cannot refresh service discovery")
					}
				}
			}
		}
	}()

	return nil
}

func (d *discovery) refresh() error {
	log.Debug("Refreshing service in discovery")

	path := d.genPathInstance()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err := d.kapi.Set(ctx, path, "", &etcd2.SetOptions{
		TTL:       ttl,
		Dir:       true,
		PrevExist: etcd2.PrevExist,
	})
	defer cancel()
	if err != nil {
		return err
	}

	return nil
}

func (d *discovery) Deregister() {

	log.Info("Deregistering service from etcd")

	path := d.genPathInstance()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err := d.kapi.Delete(ctx, path, &etcd2.DeleteOptions{
		Dir:       true,
		Recursive: true,
	})
	defer cancel()
	if err != nil {
		log.Warn("Could not deregister service", err)
	}
}

func (d *discovery) Discover(name, env, version string) (string, error) {

	log.Debugf("Discovering service %s|%s|%s", name, env, version)

	instances, err := d.list(name, env, version)
	if err != nil {
		return "", NewDiscoverError(name, env, version, err.Error())
	} else if len(instances) <= 0 {
		return "", NewDiscoverError(name, env, version, "No instances registered")
	}

	idx := rand.Intn(len(instances))

	path := instances[idx] + "/url"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := d.kapi.Get(ctx, path, nil)
	defer cancel()
	if err != nil {
		return "", NewDiscoverError(name, env, version, err.Error())
	}

	log.Debugf("Discovered service %s|%s|%s: url %s", name, env, version, resp.Node.Value)

	return resp.Node.Value, nil
}

func (d *discovery) list(name, env, version string) ([]string, error) {

	dir := fmt.Sprintf("/environments/%s/services/%s/%s/instances", env, name, version)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := d.kapi.Get(ctx, dir, nil)
	defer cancel()

	if etcd2.IsKeyNotFound(err) {
		log.Warnf("No service %s|%s|%s registered", name, env, version)
		return []string{}, nil
	} else if err != nil {
		return []string{}, err
	}

	keys := []string{}

	for _, node := range resp.Node.Nodes {
		keys = append(keys, node.Key)
	}

	return keys, nil
}

func (d *discovery) genPath() string {
	env, err := d.config.Get("env")
	if err != nil {
		log.Panicln(err)
	}

	name, err := d.config.Get("name")
	if err != nil {
		log.Panicln(err)
	}

	version, err := d.config.Get("version")
	if err != nil {
		log.Panicln(err)
	}

	return fmt.Sprintf("/environments/%s/services/%s/%s/instances", env, name, version)
}

func (d *discovery) genPathInstance() string {
	return fmt.Sprintf("%s/%s", d.genPath(), d.instanceID)
}

func (d *discovery) genPathInstanceURL() string {
	return fmt.Sprintf("%s/url", d.genPathInstance())
}
