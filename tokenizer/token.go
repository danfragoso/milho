package tokenizer

import "fmt"

type TokenType int

func (t TokenType) String() string {
	return [...]string{"Invalid", "Whitespace", "Number", "Symbol", "Reserved", "OpenParenthesis", "CloseParenthesis", "SingleQuote"}[t]
}

const (
	Invalid TokenType = iota
	Whitespace

	Number
	Symbol
	Reserved

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
	}

	if isReserved(rawToken) {
		return Reserved
	}

	for _, char := range rawToken {
		if !isDigit(char) {
			return Symbol
		}
	}

	return Number
}
