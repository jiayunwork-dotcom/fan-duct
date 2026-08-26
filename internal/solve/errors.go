package solve

import "fmt"

type SolveError struct {
	Reason string
}

func (e *SolveError) Error() string {
	return "solve: " + e.Reason
}

func IsSolveError(err error) bool {
	_, ok := err.(*SolveError)
	return ok
}

func noIntersection(why string) error {
	return &SolveError{Reason: fmt.Sprintf("no operating point: %s", why)}
}

type ParseError struct {
	Reason string
}

func (e *ParseError) Error() string {
	return "solve: " + e.Reason
}

func IsParseError(err error) bool {
	_, ok := err.(*ParseError)
	return ok
}
