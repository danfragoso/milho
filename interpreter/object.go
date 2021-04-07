package interpreter

type ObjectType int

const (
	VariableObject ObjectType = iota
	FunctionObject
)

func (r ObjectType) String() string {
	return [...]string{"Variable", "Function"}[r]
}

type Object interface {
	Type() ObjectType
	Identifier() string
	Result() Result
}

// Variable Object
type VariableObj struct {
	objectType ObjectType
	identifier string
	result     Result
}

func (v *VariableObj) Type() ObjectType {
	return v.objectType
}

func (v *VariableObj) Identifier() string {
	return v.identifier
}

func (v *VariableObj) Result() Result {
	return v.result
}

// Function Object
type FunctionObj struct {
	objectType ObjectType
	identifier string
}

func (f *FunctionObj) Type() ObjectType {
	return f.objectType
}

func (f *FunctionObj) Identifier() string {
	return f.identifier
}
