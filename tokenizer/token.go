package tokenizer

import (
	"fmt"
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
	tokenType := resolveTokenType(rawToken)
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

func resolveTokenType(rawToken string) TokenType {
	switch rawToken[0] {
	case '(':
		return OParen
	case ')':
		return CParen
	case '\'':
		return SQuote
	case '"':
		return String
	}

	if isBoolean(rawToken) {
		return Boolean
	}

	for _, char := range rawToken {
		if !isDigit(char) {
			return Symbol
		}
	}

	return Number
}
