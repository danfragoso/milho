package interpreter

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

var MILHO_STD_PATH = "./core/"

func __import(params []Expression, session *Session) (Expression, error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("import: expected at least 1 parameter, got %d, parameters must be the path to a .milho file or a symbol for a relative module or installed at MILHO_STD_PATH", len(params))
	}

	path := params[0]
	pathStr := ""

	if path.Type() == StringExpr {
		pathStr = strings.Trim(path.Value(), "\"")
	} else if path.Type() == SymbolExpr {
		pathStr = MILHO_STD_PATH + path.Value() + ".milho"
	} else {
		return nil, fmt.Errorf("import: first param to import must be either a symbol or a string got: %s", path.Value())
	}

	src, err := ioutil.ReadFile(pathStr)
	if err != nil {
		return nil, fmt.Errorf("import: failed to import module %s: %s", path.Value(), err)
	}

	tokens, err := tokenizer.Tokenize(string(src))
	if err != nil {
		return nil, fmt.Errorf("import: tokenization failed for module %s: %s", path.Value(), err)
	}

	nodes, err := parser.Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("import: parsing failed for module %s: %s", path.Value(), err)
	}

	_, err = RunFromSession(nodes, session)
	if err != nil {
		return nil, fmt.Errorf("import: session error for module %s: %s", path.Value(), err)
	}

	return createNilExpression()
}
