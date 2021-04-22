package interpreter

type Object struct {
	value      Expression
	identifier string
}

func (o *Object) Identifier() string {
	return o.identifier
}

func (o *Object) Value() Expression {
	return o.value
}
