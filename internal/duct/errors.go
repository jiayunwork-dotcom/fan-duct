package duct

import "fmt"

type ConfigError struct {
	Field  string
	Value  float64
	Wanted string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("duct: %s must satisfy %s (got %v)", e.Field, e.Wanted, e.Value)
}

func invalid(field string, value float64, wanted string) error {
	return &ConfigError{Field: field, Value: value, Wanted: wanted}
}

func IsConfigError(err error) bool {
	_, ok := err.(*ConfigError)
	return ok
}
