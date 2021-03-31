package interpreter

import "fmt"

type ResultType int

func (r ResultType) String() string {
	return [...]string{"Nil", "Number", "Boolean", "Function", "Macro", "Identifier", "List"}[r]
}

const (
	Nil ResultType = iota
	Number
	Boolean
	Function
	Macro
	Identifier
	List
)

type Result struct {
	Type  ResultType
	Value string
}

func (r *Result) String() string {
	return fmt.Sprintf("\n{Type: %s; Value: '%s'}", r.Type.String(), r.Value)
}
