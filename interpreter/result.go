package interpreter

import (
	"fmt"
	"strconv"

	"github.com/danfragoso/milho/parser"
)

type ResultType int

func (r ResultType) String() string {
	return [...]string{"Nil", "Number", "Boolean", "Function", "Macro", "Identifier", "List", "Obj", "Pending"}[r]
}

const (
	Nil ResultType = iota
	Number
	Boolean
	Function
	Macro
	Identifier
	List
	Obj
	Pending
)

type Result interface {
	Value() string
	Type() ResultType
}

func createNumberResult(value string) (*NumberResult, error) {
	parsedInt, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &NumberResult{
		value:       value,
		resultType:  Number,
		Numerator:   parsedInt,
		Denominator: 1,
	}, nil
}

type NumberResult struct {
	value      string
	resultType ResultType

	Numerator   int64
	Denominator int64
}

func (r *NumberResult) Value() string {
	return r.value
}

func (r *NumberResult) Type() ResultType {
	return r.resultType
}

func (r *NumberResult) String() string {
	return fmt.Sprintf("\n{Type: %s; Value: '%s'}", r.resultType.String(), r.value)
}

// BooleanResult

func createBooleanResult(value string) (*BooleanResult, error) {
	if value == "True" {
		return &BooleanResult{
			value:      "True",
			resultType: Boolean,
			Boolean:    true,
		}, nil
	} else if value == "False" {
		return &BooleanResult{
			value:      "False",
			resultType: Boolean,
			Boolean:    false,
		}, nil
	}

	return nil, fmt.Errorf("Wrong value %s for type boolean", value)
}

type BooleanResult struct {
	value      string
	resultType ResultType

	Boolean bool
}

func (r *BooleanResult) Value() string {
	return r.value
}

func (r *BooleanResult) Type() ResultType {
	return r.resultType
}

func (r *BooleanResult) String() string {
	return fmt.Sprintf("\n{Type: %s; Value: '%s'}", r.resultType.String(), r.value)
}

// Nil Result

func createNilResult() (*NilResult, error) {
	return &NilResult{}, nil
}

type NilResult struct{}

func (r *NilResult) Value() string {
	return "Nil"
}

func (r *NilResult) Type() ResultType {
	return Nil
}

func (r *NilResult) String() string {
	return fmt.Sprintf("\n{Type: Nil; Value: Nil}")
}

// Nil Result

func createObjectResult(obj Object) (*ObjectResult, error) {
	return &ObjectResult{
		obj,
	}, nil
}

type ObjectResult struct {
	obj Object
}

func (r *ObjectResult) Value() string {
	return r.obj.Identifier() + "#" +
		r.obj.Type().String() + ":" +
		r.obj.Result().Type().String() +
		"[" + r.obj.Result().Value() + "]"

}

func (r *ObjectResult) Type() ResultType {
	return Obj
}

func (r *ObjectResult) String() string {
	return r.Value()
}

// Pending Result
type PendingResult struct {
	Tree *parser.Node
}

func (r *PendingResult) Value() string {
	return ""
}

func (r *PendingResult) Type() ResultType {
	return Pending
}

func (r *PendingResult) String() string {
	return fmt.Sprintf("\n{Type: Pending; Value: '%s'}", r.Tree)
}

func createPendingResult(node *parser.Node) (Result, error) {
	return &PendingResult{
		Tree: node,
	}, nil
}

func createTypedResult(t ResultType, v string) (Result, error) {
	switch t {
	case Number:
		return createNumberResult(v)
	case Boolean:
		return createBooleanResult(v)
	case Nil:
		return createNilResult()
	}

	return nil, fmt.Errorf("found unresolved %s '%s'", t, v)
}
