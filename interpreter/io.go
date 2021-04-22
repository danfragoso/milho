package interpreter

import "fmt"

func __pr(params []Expression, session *Session) (Expression, error) {
	for _, param := range params {
		fmt.Print(param.Value())
	}

	return createNilExpression()
}

func __prn(params []Expression, session *Session) (Expression, error) {
	for _, param := range params {
		fmt.Print(param.Value())
	}

	fmt.Println("")
	return createNilExpression()
}
