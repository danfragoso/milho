package interpreter

var BuiltIns = map[string]*BuiltInExpression{}

func init() {
	BuiltIns = map[string]*BuiltInExpression{
		".__def":   {"Def", __def},
		".__quote": {"Quote", __quote},

		".__add": {"Add", __add},
		".__mul": {"Mul", __mul},
		".__sub": {"Sub", __sub},
		".__div": {"Div", __div},

		".__eq": {"Eq", __eq},
		".__if": {"If", __if},

		".__pr":      {"Pr", __pr},
		".__prn":     {"Prn", __prn},
		".__print":   {"Print", __print},
		".__println": {"Println", __println},

		".__str": {"Str", __str},
	}
}

var builtinInjector = `
	(.__def def .__def)
	(def quote .__quote)

	(def + .__add) (def * .__mul) (def - .__sub) (def / .__div)
	
	(def = .__eq) (def if .__if)
	
	(def pr .__pr) (def prn .__prn) (def print .__print) (def println .__println)
	
	(def str .__str)

	(def Real True)
	(def Feike False)
`
