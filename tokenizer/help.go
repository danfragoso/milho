package tokenizer

func isWhiteSpace(r rune) bool {
	switch r {
	case ' ', '\t', '\n', '\f', '\r', 0:
		return true
	}

	return false
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
