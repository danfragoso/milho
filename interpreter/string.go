package interpreter

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/danfragoso/milho/mir"
)

func __str(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	var err error
	var resultStr string
	for _, param := range params {
		param, err = evaluate(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == mir.StringExpr {
			resultStr += param.(*mir.StringExpression).Val
		} else {
			resultStr += param.Value()
		}
	}

	return mir.CreateStringExpression(resultStr)
}

func __split(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("split: expected 2 parameter, got %d, parameters must be string and separator", len(params))
	}

	str, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != mir.StringExpr {
		return nil, fmt.Errorf("split: expected first parameter to be a string, got %s", str.Type())
	}

	sep, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if str.Type() != mir.StringExpr {
		return nil, fmt.Errorf("split: expected second parameter to be a string, got %s", str.Type())
	}

	exprs := []mir.Expression{}
	strValue := strings.Trim(str.Value(), "\"")
	sepValue := strings.Trim(sep.Value(), "\"")

	strs := strings.Split(strValue, sepValue)
	for _, st := range strs {
		s, err := mir.CreateStringExpression(st)
		if err != nil {
			return nil, fmt.Errorf("split: error creating string expression from %s: %s", st, err)
		}

		exprs = append(exprs, s)
	}

	return mir.CreateListExpression(exprs...)
}

func __join(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("join: expected 2 parameter, got %d, parameters must be list of strings and separator", len(params))
	}

	lst, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if lst.Type() != mir.ListExpr {
		return nil, fmt.Errorf("join: expected first parameter to be a list, got %s", lst.Type())
	}

	var strList []string
	for _, expr := range lst.(*mir.ListExpression).Expressions {
		if expr.Type() != mir.StringExpr {
			return nil, fmt.Errorf("join: expected list to contain only strings, got %s", expr.Type())
		}

		strList = append(strList, expr.(*mir.StringExpression).Val)
	}

	sep, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if sep.Type() != mir.StringExpr {
		return nil, fmt.Errorf("join: expected second parameter to be a string, got %s", sep.Type())
	}

	sepValue := strings.Trim(sep.Value(), "\"")
	return mir.CreateStringExpression(strings.Join(strList, sepValue))
}

var randSource = rand.NewSource(time.Now().UnixNano())
var randGen = rand.New(randSource)

const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func __randomString(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("random-string: expected 1 parameter (length), got %d", len(params))
	}

	lengthExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if lengthExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("random-string: length must be a number, got %s", lengthExpr.Type())
	}

	length := int(lengthExpr.(*mir.NumberExpression).Numerator)
	if length < 0 {
		return nil, fmt.Errorf("random-string: length must be positive, got %d", length)
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphanumeric[randGen.Intn(len(alphanumeric))]
	}

	return mir.CreateStringExpression(string(result))
}

func __timestamp(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 0 {
		return nil, fmt.Errorf("timestamp: expected 0 parameters, got %d", len(params))
	}

	return mir.CreateNumberExpression(time.Now().Unix(), 1)
}

// __substring extracts a portion of a string
// substring(str:String, start:Number) -> String (from start to end)
// substring(str:String, start:Number, end:Number) -> String (from start to end)
func __substring(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) < 2 || len(params) > 3 {
		return nil, fmt.Errorf("substring: expected 2 or 3 parameters (string, start, [end]), got %d", len(params))
	}

	strExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if strExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("substring: first parameter must be a string, got %s", strExpr.Type())
	}

	startExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if startExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("substring: second parameter must be a number, got %s", startExpr.Type())
	}

	str := strExpr.(*mir.StringExpression).Val
	start := int(startExpr.(*mir.NumberExpression).Numerator)

	// Handle negative indices
	if start < 0 {
		start = 0
	}
	if start > len(str) {
		start = len(str)
	}

	// If no end parameter, return from start to end of string
	if len(params) == 2 {
		return mir.CreateStringExpression(str[start:])
	}

	// With end parameter
	endExpr, err := evaluate(params[2], session)
	if err != nil {
		return nil, err
	}

	if endExpr.Type() != mir.NumberExpr {
		return nil, fmt.Errorf("substring: third parameter must be a number, got %s", endExpr.Type())
	}

	end := int(endExpr.(*mir.NumberExpression).Numerator)
	if end < 0 {
		end = 0
	}
	if end > len(str) {
		end = len(str)
	}
	if end < start {
		end = start
	}

	return mir.CreateStringExpression(str[start:end])
}

// __replace replaces all occurrences of a substring with another
// replace(str:String, old:String, new:String) -> String
func __replace(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 3 {
		return nil, fmt.Errorf("replace: expected 3 parameters (string, old, new), got %d", len(params))
	}

	strExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if strExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("replace: first parameter must be a string, got %s", strExpr.Type())
	}

	oldExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if oldExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("replace: second parameter must be a string, got %s", oldExpr.Type())
	}

	newExpr, err := evaluate(params[2], session)
	if err != nil {
		return nil, err
	}

	if newExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("replace: third parameter must be a string, got %s", newExpr.Type())
	}

	str := strExpr.(*mir.StringExpression).Val
	old := oldExpr.(*mir.StringExpression).Val
	new := newExpr.(*mir.StringExpression).Val

	result := strings.ReplaceAll(str, old, new)
	return mir.CreateStringExpression(result)
}

// __indexOf finds the position of a substring in a string
// index-of(str:String, substr:String) -> Number (returns -1 if not found)
func __indexOf(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("index-of: expected 2 parameters (string, substring), got %d", len(params))
	}

	strExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if strExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("index-of: first parameter must be a string, got %s", strExpr.Type())
	}

	substrExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if substrExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("index-of: second parameter must be a string, got %s", substrExpr.Type())
	}

	str := strExpr.(*mir.StringExpression).Val
	substr := substrExpr.(*mir.StringExpression).Val

	index := strings.Index(str, substr)
	return mir.CreateNumberExpression(int64(index), 1)
}

// __contains checks if a string contains a substring
// contains(str:String, substr:String) -> Boolean
func __contains(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("contains: expected 2 parameters (string, substring), got %d", len(params))
	}

	strExpr, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if strExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("contains: first parameter must be a string, got %s", strExpr.Type())
	}

	substrExpr, err := evaluate(params[1], session)
	if err != nil {
		return nil, err
	}

	if substrExpr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("contains: second parameter must be a string, got %s", substrExpr.Type())
	}

	str := strExpr.(*mir.StringExpression).Val
	substr := substrExpr.(*mir.StringExpression).Val

	return mir.CreateBooleanExpression(strings.Contains(str, substr))
}
