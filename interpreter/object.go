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

func findObject(objs []*Object, identifier string) *Object {
	for _, obj := range objs {
		if obj.identifier == identifier {
			return obj
		}
	}

	return nil
}
