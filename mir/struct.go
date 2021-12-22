package mir

type Struct struct {
	Type  string
	Value map[string]Expression
}

func (*Struct) String() string {
	return "struct"
}
