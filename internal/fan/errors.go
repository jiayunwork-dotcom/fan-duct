package fan

import "fmt"

type ConfigError struct {
	Field string
	Value string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("fan: invalid %s (%s)", e.Field, e.Value)
}

func configErr(field, why string) error {
	return &ConfigError{Field: field, Value: why}
}

func IsConfigError(err error) bool {
	_, ok := err.(*ConfigError)
	return ok
}

type OutOfRangeError struct {
	Flow   float64
	Lo, Hi float64
}

func (e *OutOfRangeError) Error() string {
	return fmt.Sprintf("fan: flow %.6g m3/s is outside sample range [%.6g, %.6g] m3/s (extrapolation disabled)",
		e.Flow, e.Lo, e.Hi)
}

func IsOutOfRange(err error) bool {
	_, ok := err.(*OutOfRangeError)
	return ok
}
