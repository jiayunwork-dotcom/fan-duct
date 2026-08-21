package duct

import "fmt"

// ConfigError 报告风管配置参数不合法。字段与数值可直接贴到 CLI stderr。
type ConfigError struct {
	Field  string
	Value  float64
	Wanted string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("duct: %s must satisfy %s (got %v)", e.Field, e.Wanted, e.Value)
}

// invalid 构造一条参数校验错误。
func invalid(field string, value float64, wanted string) error {
	return &ConfigError{Field: field, Value: value, Wanted: wanted}
}

// IsConfigError 判断错误是否为参数校验错误。
func IsConfigError(err error) bool {
	_, ok := err.(*ConfigError)
	return ok
}
