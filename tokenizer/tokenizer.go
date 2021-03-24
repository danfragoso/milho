package tokenizer

import (
	"fmt"
)

func Tokenize(rawString string) ([]*Token, error) {
	s := createSource(rawString)

	var currentChar rune = s.CurrChar()
	var tokenBuffer string
	var tokens []*Token

	for currentChar != NULL_CHAR {
		if isParenthesis(currentChar) {
			tokenBuffer = string(currentChar)
		} else {
			for !isWhiteSpace(currentChar) {
				tokenBuffer += string(currentChar)

				if isParenthesis(s.PeekNextChar()) {
					break
				}

				currentChar = s.NextChar()
			}
		}

		if tokenBuffer != "" {
			currentToken, err := generateToken(tokenBuffer)
			if err != nil {
				//lint:ignore ST1005 I like my errors to be capitalized
				return nil, fmt.Errorf("Tokenization error: %s at index %d", err.Error(), s.Index)
			}

			tokens = append(tokens, currentToken)
		}

		tokenBuffer = ""
		currentChar = s.NextChar()
	}

	return tokens, nil
}