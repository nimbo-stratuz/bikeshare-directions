package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	etcd2 "go.etcd.io/etcd/client"
)

type etcd2Config struct {
	kapi etcd2.KeysAPI

	watchedMutex *sync.Mutex
	watched      map[string]string
}

func NewEtcd2Config(prefix string, conf etcd2.Config) (WritableConfig, error) {

	prefix = strings.TrimRight(prefix, "/") + "/"

	var err error
	etcdClient, err := etcd2.New(conf)
	if err != nil {
		return nil, err
	}

	keysAPI := etcd2.NewKeysAPIWithPrefix(etcdClient, prefix)

	return &etcd2Config{
		kapi: keysAPI,

		watchedMutex: &sync.Mutex{},
		watched:      make(map[string]string),
	}, nil
}

// Close closes the etcd client
func (ec *etcd2Config) Close() error {
	return nil
}

// Put a new key value pair into etcd instance
// Not accessible (etcdConfig is not exported)
func (ec *etcd2Config) Put(k string, v interface{}) (interface{}, error) {

	err := ec.setEtcd(v, k)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Get returns a string for the specified key
func (ec *etcd2Config) Get(k ...string) (string, error) {

	stringValue, err := ec.getEtcd(k...)
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *etcd2Config) GetInt(k ...string) (int, error) {

	stringValue, err := ec.getEtcd(k...)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(stringValue, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(intValue), nil
}

func (ec *etcd2Config) setEtcd(value interface{}, key ...string) error {

	fullKey := genKey(key...)

	ec.watchedMutex.Lock()

	resp, err := ec.kapi.Set(context.Background(), fullKey, fmt.Sprint(value), nil)
	if err != nil {
		ec.watchedMutex.Unlock()
		return err
	} else {
		// print common key info
		log.Printf("Set is done. Metadata is %q\n", resp)
	}

	ec.watched[fullKey] = resp.Node.Value
	ec.watchedMutex.Unlock()

	return nil
}

func (ec *etcd2Config) getEtcd(key ...string) (string, error) {

	fullKey := genKey(key...)

	// Check if the key is already watched
	ec.watchedMutex.Lock()
	val, ok := ec.watched[fullKey]
	ec.watchedMutex.Unlock()

	if ok {
		return val, nil
	}

	watcher := ec.kapi.Watcher(fullKey, nil)

	// Wait for initial value
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := watcher.Next(context.Background())
	// cancel()
	if err != nil {
		return "", err
	}

	initialValue := resp.Node.Value

	ec.watchedMutex.Lock()
	ec.watched[fullKey] = initialValue
	ec.watchedMutex.Unlock()

	// Wait for changes and update ec.watched
	go func() {
		for {
			resp, err := watcher.Next(context.Background())
			if err != nil {
				log.Printf("etcd2 Watch: %s\n", err.Error())
			}

			log.Printf("[Change: %s] Key: '%s' | Value: %s",
				resp.Action, resp.Node.Key, resp.Node.Value)

			ec.watchedMutex.Lock()
			ec.watched[fullKey] = resp.Node.Value
			ec.watchedMutex.Unlock()
		}
	}()

	return initialValue, nil
}

func genKey(key ...string) string {
	return "/" + strings.Join(key, "/")
}
