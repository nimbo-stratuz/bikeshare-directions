package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	etcd2 "go.etcd.io/etcd/client"
)

type etcd2Config struct {
	kapi etcd2.KeysAPI

	watchedMutex *sync.Mutex
	watched      map[string]watch
}

type watch struct {
	value string
	quit  chan bool
}

func NewEtcd2Config(conf etcd2.Config) (WritableConfig, error) {

	var err error
	etcdClient, err := etcd2.New(conf)
	if err != nil {
		return nil, err
	}

	keysAPI := etcd2.NewKeysAPI(etcdClient)

	return &etcd2Config{
		kapi: keysAPI,

		watchedMutex: &sync.Mutex{},
		watched:      make(map[string]watch),
	}, nil
}

func (ec *etcd2Config) getWatched(key string) (string, bool) {
	ec.watchedMutex.Lock()
	defer ec.watchedMutex.Unlock()

	if val, ok := ec.watched[key]; ok {
		return val.value, true
	}

	return "", false
}

func (ec *etcd2Config) setWatched(key string, value string) {
	ec.watchedMutex.Lock()
	defer ec.watchedMutex.Unlock()

	wtch := ec.watched[key]
	wtch = watch{
		value: value,
		quit:  wtch.quit,
	}
	ec.watched[key] = wtch
}

// Close closes the etcd client and stops all watches
func (ec *etcd2Config) Close() error {

	for _, wtch := range ec.watched {
		wtch.quit <- true
	}

	return nil
}

// Put a new key value pair into etcd instance
// Not accessible (etcdConfig is not exported)
func (ec *etcd2Config) Put(k string, v interface{}) (interface{}, error) {

	err := ec.setEtcd(k, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Get returns a string for the specified key
func (ec *etcd2Config) Get(k ...string) (string, error) {

	stringValue, err := ec.getEtcd(genKey(k...))
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *etcd2Config) GetInt(k ...string) (int, error) {

	stringValue, err := ec.getEtcd(genKey(k...))
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(stringValue, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(intValue), nil
}

func (ec *etcd2Config) setEtcd(key string, value interface{}) error {

	_, err := ec.kapi.Set(context.Background(), key, fmt.Sprint(value), nil)
	if err != nil {
		return err
	}

	// Call getEtcd to start a watch
	_, err = ec.getEtcd(key)
	if err != nil {
		return err
	}
	return nil
}

func (ec *etcd2Config) getEtcd(key string) (string, error) {

	// Check if the key is already watched
	if val, ok := ec.getWatched(key); ok {
		return val, nil
	}

	// Get initial value
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := ec.kapi.Get(ctx, key, nil)
	cancel()
	if err != nil {
		return "", err
	}

	initialValue := resp.Node.Value

	ec.watchedMutex.Lock()
	wtch := watch{
		value: initialValue,
		quit:  make(chan bool),
	}
	ec.watched[key] = wtch
	ec.watchedMutex.Unlock()

	// Wait for changes and update ec.watched
	go func() {
		log.Println("Watching key " + key)
		watcher := ec.kapi.Watcher(key, nil)
		for {
			// Create a context that will also timeout on <-quit
			ctx := context.TODO()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			go func() {
				select {
				case <-ctx.Done():
				case <-wtch.quit:
					cancel()
				}
			}()

			// Wait for a change with the defined context
			resp, err := watcher.Next(ctx)
			if err != nil {
				if err == context.Canceled {
					log.Printf("Canceled watch for key %s\n", key)
					return
				} else {
					log.Printf("etcd2 Watch: %s\n", err.Error())
				}
			}

			log.Printf("[Change: %s] Key: '%s' | Value: %s",
				resp.Action, resp.Node.Key, resp.Node.Value)

			ec.setWatched(key, resp.Node.Value)
		}
	}()

	return initialValue, nil
}

func genKey(key ...string) string {
	return strings.Join(key, "/")
}
