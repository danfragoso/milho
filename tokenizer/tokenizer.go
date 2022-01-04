package tokenizer

import (
	"errors"
	"fmt"
	"strings"
)

var line uint64 = 1
var column uint64 = 1

func Tokenize(rawString string) ([]*Token, error) {
	s := createSource(rawString)

	var currentChar rune = s.CurrChar()
	var tokenBuffer string
	var tokens []*Token

	var openString bool
	var openComment bool

	line = 1
	column = 1

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

				setPosition(currentChar, len(tokenBuffer))
				currentChar = s.NextChar()
			}
		}

		if tokenBuffer != "" {
			currentToken, err := generateToken(tokenBuffer, line, column)
			if err != nil {
				lineMarker := fmt.Sprintf(" %d:", line)

				errStr := fmt.Sprintf("Syntax error: %s \n", err.Error())
				errStr += fmt.Sprintf("%s %s\n", lineMarker, strings.Split(rawString, "\n")[line-1])
				errStr += buildIndicator(tokenBuffer, int(column-1)+len(lineMarker)-1)
				errStr += fmt.Sprintf("Line: %d, Column: %d", line, column)
				return nil, errors.New(errStr)
			}

			tokens = append(tokens, currentToken)
		}

		tokenBuffer = ""

		setPosition(currentChar, len(tokenBuffer))
		currentChar = s.NextChar()
	}

	return tokens, nil
}

func buildIndicator(indicator string, offset int) string {
	indicator = strings.Repeat("^", len(indicator))
	return strings.Repeat(" ", offset) + indicator + "\n"
}

func setPosition(currentChar rune, offset int) {
	if currentChar == NEWLINE_CHAR {
		line++
		column = 1
	} else {
		column++
	}
}
