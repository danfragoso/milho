package interpreter

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func __exit(params []Expression, session * Session) (Expression, error) {
	code := int64(0)
	if len(params) == 1 {
		if params[0].Type() != NumberExpr {
			return nil, fmt.Errorf("optional parameter code must be of type Number")
		}

		code = params[0].(*NumberExpression).Denominator
	}

	os.Exit(int(code))
	return createNilExpression()
}

func __import(params []Expression, session *Session) (Expression, error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("expected at least 1 parameter, got %d, parameters must be the path to a .milho file or a symbol for a relative module or installed at MILHO_STD_PATH", len(params))
	}

	path := params[0]
	pathStr := ""

	if path.Type() == StringExpr {
		pathStr = strings.Trim(path.Value(), "\"")
	} else if path.Type() == SymbolExpr {
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

	return createNilExpression()
}
