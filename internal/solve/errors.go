package solve

import "fmt"

// SolveError 报告工作点求解失败：无交点、超出风机样本范围、求根不收敛等。
type SolveError struct {
	Reason string
}

func (e *SolveError) Error() string {
	return "solve: " + e.Reason
}

// IsSolveError 判断错误是否为求解失败。
func IsSolveError(err error) bool {
	_, ok := err.(*SolveError)
	return ok
}

// noIntersection 构造"无工作点"错误。
func noIntersection(why string) error {
	return &SolveError{Reason: fmt.Sprintf("no operating point: %s", why)}
}

// ParseError 报告 JSON 解析失败（语法、类型、未知字段、缺字段）。
type ParseError struct {
	Reason string
}

func (e *ParseError) Error() string {
	return "solve: " + e.Reason
}

// IsParseError 判断错误是否为输入解析失败。
func IsParseError(err error) bool {
	_, ok := err.(*ParseError)
	return ok
}
