package tokenizer

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenType int

func (t TokenType) String() string {
	return [...]string{"Invalid", "Whitespace", "Number", "Symbol", "Boolean", "String", "Byte", "OpenParenthesis", "CloseParenthesis", "SingleQuote", "Comment"}[t]
}

const (
	Invalid TokenType = iota
	Whitespace

	Number
	Symbol
	Boolean
	String
	Byte

	OParen
	CParen

	SQuote
	Comment
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
		rawToken = strings.ReplaceAll(rawToken[1:len(rawToken)-1], "\\\"", "\"")
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
	case ';':
		return Comment, nil
	}

	if isDigit(rune(rawToken[0])) {
		if len(rawToken) >= 3 && rawToken[0:2] == "0x" {
			if len(rawToken) > 4 {
				return Invalid, fmt.Errorf("Invalid range for byte '%s' only allowed from 0x00 to 0xFF", rawToken)
			}

			if isHexDigit(rune(rawToken[2])) {
				if len(rawToken) == 4 && !isHexDigit(rune(rawToken[3])) {
					return Invalid, fmt.Errorf("Invalid value for byte '%s'", rawToken)
				}

				return Byte, nil
			}

			return Invalid, fmt.Errorf("Invalid value for byte '%s'", rawToken)
		}

		sArr := strings.Split(rawToken, "/")
		if len(sArr) > 2 {
			return Invalid, fmt.Errorf("Numbers can only have a single slash.")
		}

		_, err := strconv.ParseInt(sArr[0], 10, 64)
		if err != nil {
			return Invalid, fmt.Errorf("Error parsing number '%s'", sArr[0])
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
