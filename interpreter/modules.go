package interpreter

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/danfragoso/milho/parser"
	"github.com/danfragoso/milho/tokenizer"
)

func __import(params []Expression, session *Session) (Expression, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("import: expected 1 parameter, got %d, parameters must be the path to a .milho file", len(params))
	}

	path, err := evaluate(params[0], session)
	if err != nil {
		return nil, err
	}

	if path.Type() != StringExpr {
		return nil, fmt.Errorf("import: expected path to be a string, got %s", path.Type())
	}

	pathStr := strings.Trim(path.Value(), "\"")
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
