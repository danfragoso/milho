package interpreter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/danfragoso/milho/mir"
	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

var MILHO_STD_PATH = "/opt/milho/core/"

func milhoPath() string {
	envPath := os.Getenv("MILHO_STD_PATH")
	if envPath != "" {
		return envPath
	}

	fmt.Printf("Warning, MILHO_STD_PATH is not defined, using default location %s\n", MILHO_STD_PATH)
	return MILHO_STD_PATH
}

func __exit(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	code := int64(0)
	if len(params) == 1 {
		if params[0].Type() != mir.NumberExpr {
			return nil, fmt.Errorf("optional parameter code must be of type Number")
		}

		code = params[0].(*mir.NumberExpression).Denominator
	}

	os.Exit(int(code))
	return mir.CreateNilExpression()
}

func __import(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("expected at least 1 parameter, got %d, parameters must be the path to a .milho file or a symbol for a relative module or installed at MILHO_STD_PATH", len(params))
	}

	path := params[0]
	pathStr := ""

	if path.Type() == mir.StringExpr {
		pathStr = strings.Trim(path.Value(), "\"")
	} else if path.Type() == mir.SymbolExpr {
		pathStr = milhoPath() + path.Value() + ".milho"
	} else {
		return nil, fmt.Errorf("first param to import must be either a symbol or a string got: %s", path.Value())
	}

	src, err := ioutil.ReadFile(pathStr)
	if err != nil {
		return nil, fmt.Errorf("failed to import module %s: %s", path.Value(), err)
	}

	tokens, err := tokenizer.Tokenize(string(src))
	if err != nil {
		return nil, fmt.Errorf("tokenization failed for module %s: %s", path.Value(), err)
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("parsing failed for module %s: %s", path.Value(), err)
	}

	_, err = RunFromSession(nodes, session)
	if err != nil {
		return nil, fmt.Errorf("session error for module %s: %s", path.Value(), err)
	}

	return mir.CreateNilExpression()
}

func __read(params []mir.Expression, session *mir.Session) (mir.Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("expected 1 parameter, got %d", len(params))
	}

	expr, err := evaluate(params[0], session)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate parameter: %s", err)
	}

	if expr.Type() != mir.StringExpr {
		return nil, fmt.Errorf("parameter must be of type string, got: %s", expr.Value())
	}

	pathStr := expr.(*mir.StringExpression).Val
	src, err := ioutil.ReadFile(pathStr)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", pathStr, err)
	}

	return mir.CreateStringExpression(string(src))
}
