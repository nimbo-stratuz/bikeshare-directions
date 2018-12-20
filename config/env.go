package config

import (
	"errors"
	"os"
	"strconv"
)

// EnvConfig is a client for reading env variables
type EnvConfig struct {
}

// NewEnvConfig New creates an EnvConfig instance
func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}

func (ec *EnvConfig) Close() error {
	return nil
}

// Get returns a string for the specified key
func (ec *EnvConfig) Get(k string) (string, error) {

	stringValue, err := ec.getEnv(k)
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *EnvConfig) GetInt(k string) (int, error) {

	stringValue, err := ec.getEnv(k)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(stringValue, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(intValue), nil
}

func (ec *EnvConfig) getEnv(key string) (string, error) {

	value := os.Getenv(key)

	if value == "" {
		return "", errors.New("key " + key + " not found")
	}

	return value, nil
}
