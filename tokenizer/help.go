package tokenizer

func isWhiteSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\f', '\r', 0:
		return true
	}

	return false
}

func isNewLine(r rune) bool {
	return r == '\n'
}

func isCommentStart(r rune) bool {
	return r == ';'
}

func isSingleQuote(r rune) bool {
	return r == '\''
}

func isDoubleQuote(r rune) bool {
	return r == '"'
}

func isBackslash(r rune) bool {
	return r == '\\'
}

func isParenthesis(r rune) bool {
	switch r {
	case '(', ')':
		return true
	}

	return false
}

func isDigit(r rune) bool {
	if r >= 48 && r <= 57 {
		return true
	}

	return false
}

func isHexDigit(r rune) bool {
	if isDigit(r) || r >= 65 && r <= 70 || r >= 97 && r <= 102 {
		return true
	}

	return false
}
