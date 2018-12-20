package config

import (
	"log"
	"errors"
	"strconv"
	"strings"
)

// Config is a common interface that ensures basic methods
type Config interface {
	Close() error
	Get(...string) (string, error)
	GetInt(...string) (int, error)
}

// MultiConfig represents a hierarchy of Configs
type MultiConfig struct {
	configs []Config
}

// New creates a new MultiConfig from the specified Configs
// When querying the config, Configs are checked for the specified key
// in the same order as they are specified here.
func New(configs ...Config) MultiConfig {
	return MultiConfig{configs: configs}
}

// Close closes all underlying Configs
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

// Get returns a string value for the specified key
func (mc *MultiConfig) Get(key ...string) (string, error) {

	for _, c := range mc.configs {
		value, err := c.Get(key...)
		if err != nil {
			// key found
			return value, nil
		}
	}

	return "", errors.New("Key " + strings.Join(key, ".") + " not found")
}

// GetInt returns an int value for the specified key
func (mc *MultiConfig) GetInt(key ...string) (int, error) {

	for _, c := range mc.configs {
		value, err := c.GetInt(key...)
		if err, ok := err.(*strconv.NumError); ok {
			// key found, but value is not an integer
			return 0, err
		}
		if err != nil {
			log.Println("MultiConfig.GetInt " + err.Error())
		}
		if err == nil {
			// key found
			return value, nil
		}
	}

	return 0, errors.New("Key " + strings.Join(key, ".") + " not found")
}
