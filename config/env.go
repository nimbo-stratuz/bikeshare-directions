package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// envConfig is a client for reading env variables
type envConfig struct {
}

// NewEnvConfig New creates an envConfig instance
func NewEnvConfig() Config {
	return &envConfig{}
}

// Close does nothing for envConfig
func (ec *envConfig) Close() error {
	return nil
}

// Get returns a string for the specified key
func (ec *envConfig) Get(key ...string) (string, error) {

	stringValue, err := ec.getEnv(key...)
	if err != nil {
		return "", err
	}

	return stringValue, nil
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *envConfig) GetInt(key ...string) (int, error) {

	stringValue, err := ec.getEnv(key...)
	if err != nil {
		return 0, err
	}

	intValue, err := strconv.ParseInt(stringValue, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(intValue), nil
}

func (ec *envConfig) getEnv(key ...string) (string, error) {

	fullKey := strings.ToUpper(strings.Join(key, "_"))

	value := os.Getenv(fullKey)

	if value == "" {
		return "", errors.New("key " + fullKey + " not found")
	}

	return value, nil
}
