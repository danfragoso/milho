package interpreter

var BuiltIns = map[string]*BuiltInExpression{}

func init() {
	BuiltIns = map[string]*BuiltInExpression{
		".__def":   {"Def", __def},
		".__defn":  {"Defn", __defn},
		".__quote": {"Quote", __quote},
		".__let":   {"Let", __let},
		".__type":  {"Type", __type},
		".__fn":    {"Fn", __fn},
		".__time":  {"Time", __time},
		".__progn": {"Progn", __progn},

		".__add": {"Add", __add},
		".__mul": {"Mul", __mul},
		".__sub": {"Sub", __sub},
		".__div": {"Div", __div},

		".__eq": {"Eq", __eq},
		".__if": {"If", __if},

		".__car": {"Car", __car},
		".__cdr": {"Cdr", __cdr},

		".__pr":      {"Pr", __pr},
		".__prn":     {"Prn", __prn},
		".__print":   {"Print", __print},
		".__println": {"Println", __println},

		".__str": {"Str", __str},
	}
}

var builtinInjector = `
	(.__def def .__def) (def defn .__defn) (def quote .__quote) (def type .__type)
	(def let .__let) (def fn .__fn) (def time .__time) (def progn .__progn)

	(def + .__add) (def * .__mul) (def - .__sub) (def / .__div)
	
	(def = .__eq) (def if .__if)
	
	(def car .__car) (def cdr .__cdr)

	(def pr .__pr) (def prn .__prn) (def print .__print) (def println .__println)
	
	(def str .__str)

	(def Real True)
	(def Feike False)

	(def Nil ())
`
var functionInjector = `
(defn test (name expected result)
	(if (= expected result)
		(println "PASS:" name)
		(progn
			(println "FAIL:" name)
			(println "` + "\u200e" + `  └─ Value {" (str result) "} doesn't equal expected result {" (str expected) "}."))))
`
