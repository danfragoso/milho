package interpreter

func __str(params []Expression, session *Session) (Expression, error) {
	var err error
	var resultStr string
	for _, param := range params {
		param, err = resolveExpression(param, session)
		if err != nil {
			return nil, err
		}

		if param.Type() == StringExpr {
			resultStr += param.(*StringExpression).Val
		} else {
			resultStr += param.Value()
		}
	}

	return createStringExpression(resultStr)
}
