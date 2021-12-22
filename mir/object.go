package mir

type Object struct {
	value      Expression
	identifier string
}

func CreateObject(value Expression, identifier string) *Object {
	return &Object{
		identifier: identifier,
		value:      value,
	}
}

func (o *Object) Identifier() string {
	return o.identifier
}

func (o *Object) Value() Expression {
	return o.value
}

func FindObject(objs []*Object, identifier string) *Object {
	for _, obj := range objs {
		if obj.identifier == identifier {
			return obj
		}
	}

	return nil
}
