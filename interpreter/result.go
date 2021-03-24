package interpreter

import "fmt"

type ResultType int

func (r ResultType) String() string {
	return [...]string{"Nil", "Number"}[r]
}

const (
	Nil ResultType = iota
	Number
)

type Result struct {
	Type  ResultType
	Value string
}

func (r *Result) String() string {
	return fmt.Sprintf("\n{Type: %s; Value: '%s'}", r.Type, r.Value)
}
