package tokenizer

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

func (t TokenType) String() string {
	return [...]string{"Invalid", "Whitespace", "Number", "Symbol", "Boolean", "String", "OpenParenthesis", "CloseParenthesis", "SingleQuote"}[t]
}

const (
	Invalid TokenType = iota
	Whitespace

	Number
	Symbol
	Boolean
	String

	OParen
	CParen

	SQuote
)

type Token struct {
	Type  TokenType
	Value string
}

func generateToken(rawToken string) (*Token, error) {
	tokenType, err := resolveTokenType(rawToken)
	if err != nil {
		return nil, err
	}

	switch tokenType {
	case Invalid:
		return nil, fmt.Errorf("invalid token '%s'", rawToken)
	case String:
		rawToken = strings.Trim(rawToken, "\"")
	}

	return &Token{
		Value: rawToken,
		Type:  tokenType,
	}, nil
}

func resolveTokenType(rawToken string) (TokenType, error) {
	switch rawToken[0] {
	case '(':
		return OParen, nil
	case ')':
		return CParen, nil
	case '\'':
		return SQuote, nil
	case '"':
		return String, nil
	}

	if isDigit(rune(rawToken[0])) {
		sArr := strings.Split(rawToken, "/")
		if len(sArr) > 2 {
			return Invalid, fmt.Errorf("Numbers can only have a single slash.")
		}

		_, err := strconv.ParseInt(sArr[0], 10, 64)
		if err != nil {
			return Invalid, fmt.Errorf("Error parsing numerator '%s'", sArr[0])
		}

		if len(sArr) == 2 {
			if sArr[1] == "" {
				return Number, nil
			}

			_, err = strconv.ParseInt(sArr[1], 10, 64)
			if err != nil {
				return Invalid, fmt.Errorf("Error parsing denominator '%s'", sArr[1])
			}
		}

		return Number, nil
	}

	if isBoolean(rawToken) {
		return Boolean, nil
	}

	return Symbol, nil
}
