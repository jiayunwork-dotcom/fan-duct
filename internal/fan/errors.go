package fan

import "fmt"

// ConfigError 报告风机曲线配置不合法。字段与数值可直接贴到 CLI stderr。
type ConfigError struct {
	Field string
	Value string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("fan: invalid %s (%s)", e.Field, e.Value)
}

// configErr 构造一条曲线校验错误。
func configErr(field, why string) error {
	return &ConfigError{Field: field, Value: why}
}

// IsConfigError 判断错误是否为风机曲线参数错误。
func IsConfigError(err error) bool {
	_, ok := err.(*ConfigError)
	return ok
}

// OutOfRangeError 报告流量超出风机样本范围且外推策略禁止外推。
type OutOfRangeError struct {
	Flow   float64
	Lo, Hi float64
}

func (e *OutOfRangeError) Error() string {
	return fmt.Sprintf("fan: flow %.6g m3/s is outside sample range [%.6g, %.6g] m3/s (extrapolation disabled)",
		e.Flow, e.Lo, e.Hi)
}

// IsOutOfRange 判断错误是否为样本范围越界。
func IsOutOfRange(err error) bool {
	_, ok := err.(*OutOfRangeError)
	return ok
}
