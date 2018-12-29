package config

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Config is a common interface that ensures basic methods
type Config interface {
	Close() error
	Get(...string) (string, error)
	GetInt(...string) (int, error)
}

// WritableConfig is an editable config source (supports method Put)
type WritableConfig interface {
	Config
	Put(string, interface{}) (interface{}, error)
}

// multiConfig represents a hierarchy of Configs
type multiConfig struct {
	configs []Config
}

// New creates a new MultiConfig from the specified Configs
// When querying the config, Configs are checked for the specified key
// in the same order as they are specified here.
func New(configs ...Config) (Config, error) {

	if len(configs) <= 0 {
		return nil, fmt.Errorf("At least one Config has to be provided")
	}

	return &multiConfig{configs: configs}, nil
}

// Close closes all underlying Configs
func (mc *multiConfig) Close() error {

	log.Info("Closing multiConfig")

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
func (mc *multiConfig) Get(key ...string) (string, error) {

	var errs []string

	for _, c := range mc.configs {
		value, err := c.Get(key...)
		if err == nil {
			// key found
			return value, nil
		}
		errs = append(errs, err.Error())
	}

	return "", fmt.Errorf("Key not found: %s | Errors: [%s]", strings.Join(key, "."), strings.Join(errs, "; "))
}

// GetInt returns an int value for the specified key
func (mc *multiConfig) GetInt(key ...string) (int, error) {

	var errs []string

	for _, c := range mc.configs {
		value, err := c.GetInt(key...)
		if err == nil {
			// key found
			return value, nil
		}
		errs = append(errs, err.Error())
	}

	return 0, fmt.Errorf("Key not found: %s | Errors: [%s]", strings.Join(key, "."), strings.Join(errs, "; "))
}
