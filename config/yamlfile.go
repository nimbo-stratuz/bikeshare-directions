package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type yamlObject map[interface{}]interface{}

// YamlFileConfig is a client for reading env variables
type YamlFileConfig struct {
	yaml yamlObject
}

// NewYamlFileConfig New creates an YamlFileConfig instance
func NewYamlFileConfig(filePath string) (*YamlFileConfig, error) {

	yamlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	yamlMap := make(yamlObject)

	if err := yaml.Unmarshal(yamlBytes, &yamlMap); err != nil {
		return nil, err
	}

	return &YamlFileConfig{
		yaml: yamlMap,
	}, nil
}

// Close does nothing for YamlFileConfig
func (ec *YamlFileConfig) Close() error {
	return nil
}

// Get returns a string for the specified key
func (ec *YamlFileConfig) Get(key ...string) (string, error) {

	value, err := ec.getYamlValue(key...)
	if err != nil {
		return "", err
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return string(v), nil
	default:
		return "", fmt.Errorf("Unsupported type")

	}
}

// GetInt returns a string for the specified key converted to a 32 bit integer
func (ec *YamlFileConfig) GetInt(key ...string) (int, error) {

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

func (ec *YamlFileConfig) getYamlValue(key ...string) (interface{}, error) {

	fullKey := strings.ToLower(strings.Join(key, "."))

	var objPtr interface{}
	objPtr = ec.yaml

	for _, k := range key {
		switch i := objPtr.(type) {
		case yamlObject:
			if val, ok := i[k]; ok {
				objPtr = val
			} else {
				break
			}
		default:
			break
		}
	}

	// Check if current objPtr is a final value (not a nested object)
	if _, notok := objPtr.(yamlObject); notok {
		return "", fmt.Errorf("Not a final value: %s", fullKey)
	}
	return objPtr, nil
}
