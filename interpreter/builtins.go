package interpreter

import "github.com/danfragoso/milho/mir"

func init() {
	mir.BuiltIns = map[string]*mir.BuiltInExpression{
		".__def":   {"Def", __def},
		".__defn":  {"Defn", __defn},
		".__quote": {"Quote", __quote},
		".__let":   {"Let", __let},
		".__type":  {"Type", __type},
		".__fn":    {"Fn", __fn},
		".__time":  {"Time", __time},
		".__do":    {"do", __do},
		".__eval":  {"Eval", __eval},

		".__add": {"Add", __add},
		".__mul": {"Mul", __mul},
		".__sub": {"Sub", __sub},
		".__div": {"Div", __div},
		".__lt":  {"Lt", __lt},

		".__eq":     {"Eq", __eq},
		".__negate": {"Negate", __negate},
		".__if":     {"If", __if},

		".__car": {"Car", __car},
		".__cdr": {"Cdr", __cdr},

		".__pr":      {"Pr", __pr},
		".__prn":     {"Prn", __prn},
		".__print":   {"Print", __print},
		".__println": {"Println", __println},

		".__list": {"List () -> Nil", __list},
		".__cons": {"Cons", __cons},
		".__map":  {"Map", __map},

		".__split": {"Split", __split},
		".__join":  {"Join (stringList:(String) joiner:String) -> String", __join},
		".__str":   {"Str (?params:...) -> String", __str},

		".__exec": {"Exec (command:String ?params:...String) -> String", __exec},

		".__exec-code":   {"Exec (command:String ?params:...String) -> String", __execCode},
		".__exec-stdout": {"Exec (command:String ?params:...String) -> String", __execStdout},
		".__exec-stderr": {"Exec (command:String ?params:...String) -> String", __execStderr},

		".__import": {"Import (module:String|Symbol ?namespace:Symbol) -> Nil", __import},

		".__sleep": {"Sleep (ms:Number) -> Nil", __sleep},

		".__range": {"Range (min:Number max:Number) -> Nil", __range},
		".__exit":  {"Exit (?code:Number) -> Nil", __exit},

		".__push": {"Push (item:Any list:List) -> List", __push},

		".__map-create":    {"MapCreate (?params:...) -> Map", __mapCreate},
		".__map-get":       {"MapGet (map:Map key:String) -> Any", __mapGet},
		".__map-set":       {"MapSet (map:Map key:String value:Any) -> Map", __mapSet},
		".__map-delete":    {"MapDelete (map:Map key:String) -> Map", __mapDelete},
		".__map-keys":      {"MapKeys (map:Map) -> List", __mapKeys},
		".__map-from-json": {"MapFromJSON (json:String) -> Map", __mapFromJSON},

		".__read":  {"Read (path:String) -> String", __read},
		".__match": {"Match (value:Any pattern:Any) -> Any", __match},
	}
}

var builtinInjector = `
	(.__def def .__def) (def defn .__defn) (def quote .__quote) (def type .__type)
	(def let .__let) (def fn .__fn) (def time .__time) (def do .__do)

	(def + .__add) (def * .__mul) (def - .__sub) (def / .__div) (def < .__lt)

	(def = .__eq) (def ! .__negate) (def if .__if) (def map .__map)

	(def car .__car) (def cdr .__cdr) (def cons .__cons)

	(def pr .__pr) (def prn .__prn) (def print .__print) (def println .__println)

	(def eval .__eval)
	(def push .__push)
	(def str .__str)
	(def exit .__exit)
	(def Real True)
	(def Feike False)

	(def map-create .__map-create)
	(def map-get .__map-get)
	(def map-set .__map-set)
	(def map-delete .__map-delete)
	(def map-keys .__map-keys)
	(def map-from-json .__map-from-json)

	(def match .__match)
	(def read .__read)
	
	(def exec-code .__exec-code)
	(def exec-stdout .__exec-stdout)
	(def exec-stderr .__exec-stderr)

	(def Nil ())

	(def list .__list)

	(def split .__split)
	(def join .__join)

	(def exec .__exec)

	(def import .__import)

	(def sleep .__sleep)
	(def range .__range)
`
var functionInjector = `
(defn Number? (n)
	(= (type n) 'Number))

(defn String? (n)
	(= (type n) 'String))

(defn Bool? (n)
	(= (type n) 'Bool))

(defn Symbol? (n)
	(= (type n) 'Symbol))

(defn fat (n)
	(if (= n 0)
		1
		(* n (fat (- n 1)))))

(defn test (name expected result)
	(if (= expected result)
		(println "PASS:" name)
		(do
			(println "FAIL:" name)
			(println "` + "\u200e" + ` └─ Value {" (str result) "} doesn't equal expected result {" (str expected) "}."))))
`
