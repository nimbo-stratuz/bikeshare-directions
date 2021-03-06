package config

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	etcd3 "go.etcd.io/etcd/clientv3"
)

// etcdConfig is a client for communication with etcd
type etcdConfig struct {
	cli    *etcd3.Client
	prefix string
}

// NewEtcdConfig creates an etcdConfig instance
func NewEtcdConfig(prefix string, conf etcd3.Config) (WritableConfig, error) {

	prefix = strings.TrimRight(prefix, "/") + "/"

	var err error
	etcdClient, err := etcd3.New(conf)
	if err != nil {
		return nil, err
	}

	return &etcdConfig{
		cli:    etcdClient,
		prefix: prefix,
	}, nil
}

// Close closes the etcd client
func (ec *etcdConfig) Close() error {
	log.Println("Closing etcdConfig")
	return ec.cli.Close()
}

// Put a new key value pair into etcd instance
// Not accessible (etcdConfig is not exported)
func (ec *etcdConfig) Put(k string, v interface{}) (interface{}, error) {

	err := ec.setEtcd(k, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Get returns a string for the specified key
func (ec *etcdConfig) Get(k ...string) (string, error) {

	stringValue, err := ec.getEtcd(k...)
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *etcdConfig) GetInt(k ...string) (int, error) {

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

func (ec *etcdConfig) setEtcd(key string, value interface{}) error {

	key = ec.prefix + strings.TrimLeft(key, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err := ec.cli.Put(ctx, key, fmt.Sprint(value))
	cancel()
	if err != nil {
		return err
	}

	return nil
}

func (ec *etcdConfig) getEtcd(key ...string) (string, error) {

	fullKey := strings.ToLower(strings.Join(key, "/"))
	fullKey = ec.prefix + strings.TrimLeft(fullKey, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := ec.cli.Get(ctx, fullKey)
	cancel()
	if err != nil {
		return "", err
	}

	for _, ev := range resp.Kvs {
		if string(ev.Key) == fullKey {
			return string(ev.Value), nil
		}
	}

	return "", errors.New("key " + fullKey + " not found")
}
