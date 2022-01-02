package tokenizer

import (
	"fmt"
)

func Tokenize(rawString string) ([]*Token, error) {
	s := createSource(rawString)

	var currentChar rune = s.CurrChar()
	var tokenBuffer string
	var tokens []*Token

	var openString bool
	var openComment bool

	for currentChar != NULL_CHAR {
		if isParenthesis(currentChar) || isSingleQuote(currentChar) {
			tokenBuffer = string(currentChar)
		} else {
			for openComment || openString || !isWhiteSpace(currentChar) {
				if isDoubleQuote(currentChar) && !isBackslash(s.PreviousChar()) {
					openString = !openString
				}

				if isCommentStart(currentChar) {
					openComment = true
				}

				if openComment && isNewLine(currentChar) {
					openComment = false
				}

				tokenBuffer += string(currentChar)

				if s.PeekNextChar() == NULL_CHAR {
					break
				}

				if isParenthesis(s.PeekNextChar()) && !openString {
					break
				}

				currentChar = s.NextChar()
			}
		}

		if tokenBuffer != "" {
			currentToken, err := generateToken(tokenBuffer)
			if err != nil {
				return nil, fmt.Errorf("Tokenization error: %s at index %d", err.Error(), s.Index)
			}

			tokens = append(tokens, currentToken)
		}

		tokenBuffer = ""
		currentChar = s.NextChar()
	}

	return tokens, nil
}
