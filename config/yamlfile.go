package config

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type yamlObject map[interface{}]interface{}

// yamlFileConfig is a client for reading env variables
type yamlFileConfig struct {
	yaml yamlObject
}

// NewYamlFileConfig New creates an yamlFileConfig instance
func NewYamlFileConfig(filePath string) (Config, error) {

	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	yamlMap := make(yamlObject)

	if err := yaml.Unmarshal(yamlBytes, &yamlMap); err != nil {
		return nil, err
	}

	return &yamlFileConfig{
		yaml: yamlMap,
	}, nil
}

// Close does nothing for yamlFileConfig
func (ec *yamlFileConfig) Close() error {
	return nil
}

// Get returns a string for the specified key
func (ec *yamlFileConfig) Get(key ...string) (string, error) {

	value, err := ec.getYamlValue(key...)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	default:
		return "", fmt.Errorf("Unsupported type")

	}
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *yamlFileConfig) GetInt(key ...string) (int, error) {

	value, err := ec.getYamlValue(key...)
	if err != nil {
		return 0, err
	}

	switch v := value.(type) {
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("Wanted int, got string")

	}
}

func (ec *yamlFileConfig) getYamlValue(key ...string) (interface{}, error) {

	fullKey := strings.ToLower(strings.Join(key, "."))

	var objPtr interface{}
	objPtr = ec.yaml

	for idx, k := range key {
		switch i := objPtr.(type) {
		case yamlObject:
			if val, ok := i[k]; ok {
				objPtr = val
			}
		default:
			return "", fmt.Errorf("Key not a YAML object: %s", strings.Join(key[0:idx], "."))
		}
	}

	// Check if current objPtr is a final value (not a nested object)
	if _, notok := objPtr.(yamlObject); notok {
		return "", fmt.Errorf("Not a final value: %s", fullKey)
	}
	return objPtr, nil
}
