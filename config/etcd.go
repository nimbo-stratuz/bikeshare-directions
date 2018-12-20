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

// EtcdConfig is a client for communication with etcd
type EtcdConfig struct {
	cli    *etcd3.Client
	prefix string
}

// New creates an EtcdConfig instance
func New(prefix string, conf etcd3.Config) (*EtcdConfig, error) {

	prefix = strings.TrimRight(prefix, "/") + "/"

	var err error
	etcdClient, err := etcd3.New(conf)
	if err != nil {
		return nil, err
	}

	return &EtcdConfig{
		cli:    etcdClient,
		prefix: prefix,
	}, nil
}

// Close closes the etcd client
func (ec *EtcdConfig) Close() error {
	log.Println("Closing EtcdConfig")
	return ec.cli.Close()
}

// Put a new key value pair into etcd instance
func (ec *EtcdConfig) Put(k string, v interface{}) (interface{}, error) {

	err := ec.setEtcd(k, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Get returns a string for the specified key
func (ec *EtcdConfig) Get(k string) (string, error) {

	stringValue, err := ec.getEtcd(k)
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *EtcdConfig) GetInt(k string) (int, error) {

	stringValue, err := ec.getEtcd(k)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(stringValue, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(intValue), nil
}

func (ec *EtcdConfig) setEtcd(key string, value interface{}) error {

	key = ec.prefix + strings.TrimLeft(key, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	_, err := ec.cli.Put(ctx, key, fmt.Sprint(value))
	cancel()
	if err != nil {
		return err
	}

	return nil
}

func (ec *EtcdConfig) getEtcd(key string) (string, error) {

	key = ec.prefix + strings.TrimLeft(key, "/")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	resp, err := ec.cli.Get(ctx, key)
	cancel()
	if err != nil {
		return "", err
	}

	for _, ev := range resp.Kvs {
		if string(ev.Key) == key {
			log.Printf("Found key %s\n", string(ev.Key))
			return string(ev.Value), nil
		}
	}

	return "", errors.New("key " + key + " not found")
}
