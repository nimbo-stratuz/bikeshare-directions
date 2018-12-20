package config

import (
	"errors"
	"log"
	"strings"
)

type Config interface {
	Close() error
	Get(string) (string, error)
	GetInt(string) (int, error)
}

type MultiConfig struct {
	configs []Config
}

func New(configs ...Config) MultiConfig {
	return MultiConfig{configs: configs}
}

func (mc *MultiConfig) Close() error {

	var errs []string

	for _, c := range mc.configs {
		if err := c.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (mc *MultiConfig) Get(key string) (string, error) {

	for _, c := range mc.configs {
		value, err := c.Get(key)
		if err == nil {
			// key found
			return value, nil
		}
	}

	return "", errors.New("Key " + key + " not found")
}

func (mc *MultiConfig) GetInt(key string) (int, error) {

	for _, c := range mc.configs {
		value, err := c.GetInt(key)
		if err == nil {
			// key found
			return value, nil
		}
	}

	return 0, errors.New("Key " + key + " not found")
}
