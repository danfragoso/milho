package tokenizer

const NULL_CHAR rune = 0

type Source struct {
	Length int
	Index  int
	Value  string
}

func createSource(src string) *Source {
	return &Source{
		Length: len(src),
		Value:  src,
	}
}

func (src *Source) CurrChar() rune {
	if src.Index < src.Length {
		return rune(src.Value[src.Index])
	}

	return NULL_CHAR
}

func (src *Source) NextChar() rune {
	src.Index++
	if src.Index < src.Length {
		return rune(src.Value[src.Index])
	}

	return NULL_CHAR
}

func (src *Source) PeekNextChar() rune {
	if src.Index+1 < src.Length {
		return rune(src.Value[src.Index+1])
	}

	return NULL_CHAR
}
