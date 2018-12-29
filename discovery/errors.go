package discovery

import "fmt"

// DiscoverError is returned when a service fails to discover
// Some other service, specified by ser
type DiscoverError struct {
	service, env, version string // The service that cannot be discovered
	reason                string
}

// NewDiscoverError creates a new DiscoverError
func NewDiscoverError(service, env, version, reason string) *DiscoverError {
	return &DiscoverError{service, env, version, reason}
}

func (de *DiscoverError) Error() string {
	return fmt.Sprintf("Cannot discover service %s|%s|%s: %s", de.service, de.env, de.version, de.reason)
}

// RegisterError is returned when a service fails to
// register with etcd.
type RegisterError struct {
	reason string
}

// NewRegisterError creates a new RegisterError
func NewRegisterError(reason string) *RegisterError {
	return &RegisterError{reason}
}

func (re *RegisterError) Error() string {
	return fmt.Sprintf("Cannot register service: %s", re.reason)
}
