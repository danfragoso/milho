package interpreter

type Struct struct {
	Type  string
	Value map[string]Expression
}

func (*Struct) String() string {
	return "struct"
}

func __struct(params []Expression, session *Session) (Expression, error) {
	return createNumberExpression(0, 1)
}
