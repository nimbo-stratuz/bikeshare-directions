package config

type Config interface {
	Put(string, interface{}) (interface{}, error)
	Get(string) (string, error)
	GetInt(string) (int, error)
}

