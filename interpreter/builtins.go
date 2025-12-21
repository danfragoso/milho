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

		".__split":         {"Split", __split},
		".__join":          {"Join (stringList:(String) joiner:String) -> String", __join},
		".__str":           {"Str (?params:...) -> String", __str},
		".__random-string": {"RandomString (length:Number) -> String", __randomString},
		".__timestamp":     {"Timestamp () -> Number", __timestamp},
		".__substring":     {"Substring (str:String start:Number ?end:Number) -> String", __substring},
		".__replace":       {"Replace (str:String old:String new:String) -> String", __replace},
		".__index-of":      {"IndexOf (str:String substr:String) -> Number", __indexOf},
		".__contains":      {"Contains (str:String substr:String) -> Boolean", __contains},

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

		".__http-server-create":      {"HttpServerCreate (port:Number) -> Number", __httpServerCreate},
		".__http-server-route":       {"HttpServerRoute (serverID:Number method:String path:String handler:Function) -> Nil", __httpServerRoute},
		".__http-server-start":       {"HttpServerStart (serverID:Number) -> Nil", __httpServerStart},
		".__http-server-start-async": {"HttpServerStartAsync (serverID:Number) -> Nil", __httpServerStartAsync},
		".__http-server-stop":        {"HttpServerStop (serverID:Number) -> Nil", __httpServerStop},

		".__http-request-method":  {"HttpRequestMethod (request:Map) -> String", __httpRequestMethod},
		".__http-request-path":    {"HttpRequestPath (request:Map) -> String", __httpRequestPath},
		".__http-request-headers": {"HttpRequestHeaders (request:Map) -> Map", __httpRequestHeaders},
		".__http-request-body":    {"HttpRequestBody (request:Map) -> String", __httpRequestBody},
		".__http-request-query":   {"HttpRequestQuery (request:Map) -> Map", __httpRequestQuery},

		".__http-response-status": {"HttpResponseStatus (response:Map status:Number) -> Map", __httpResponseStatus},
		".__http-response-header": {"HttpResponseHeader (response:Map key:String value:String) -> Map", __httpResponseHeader},
		".__http-response-body":   {"HttpResponseBody (response:Map body:String) -> Map", __httpResponseBody},
		".__http-response-json":   {"HttpResponseJson (response:Map jsonData:String|Map) -> Map", __httpResponseJson},

		".__tcp-connect":   {"TcpConnect (host:String port:Number) -> Number", __tcpConnect},
		".__tcp-write":     {"TcpWrite (socketID:Number data:String) -> Nil", __tcpWrite},
		".__tcp-read":      {"TcpRead (socketID:Number ?delimiter:String) -> String", __tcpRead},
		".__tcp-read-line": {"TcpReadLine (socketID:Number) -> String", __tcpReadLine},
		".__tcp-close":     {"TcpClose (socketID:Number) -> Nil", __tcpClose},
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
	(def random-string .__random-string)
	(def timestamp .__timestamp)
	(def substring .__substring)
	(def replace .__replace)
	(def index-of .__index-of)
	(def contains .__contains)

	(def exec .__exec)

	(def import .__import)

	(def sleep .__sleep)
	(def range .__range)

	(def http-server-create .__http-server-create)
	(def http-server-route .__http-server-route)
	(def http-server-start .__http-server-start)
	(def http-server-start-async .__http-server-start-async)
	(def http-server-stop .__http-server-stop)

	(def http-request-method .__http-request-method)
	(def http-request-path .__http-request-path)
	(def http-request-headers .__http-request-headers)
	(def http-request-body .__http-request-body)
	(def http-request-query .__http-request-query)

	(def http-response-status .__http-response-status)
	(def http-response-header .__http-response-header)
	(def http-response-body .__http-response-body)
	(def http-response-json .__http-response-json)

	(def tcp-connect .__tcp-connect)
	(def tcp-write .__tcp-write)
	(def tcp-read .__tcp-read)
	(def tcp-read-line .__tcp-read-line)
	(def tcp-close .__tcp-close)
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
