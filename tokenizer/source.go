package tokenizer

const NULL_CHAR rune = 0

type Source struct {
	Length int
	Index  int
	Value  []rune
}

func createSource(src string) *Source {
	runeArr := []rune(src)
	return &Source{
		Length: len(runeArr),
		Value:  runeArr,
	}
}

func (src *Source) CurrChar() rune {
	if src.Index < src.Length {
		return src.Value[src.Index]
	}

	return NULL_CHAR
}

func (src *Source) NextChar() rune {
	src.Index++
	if src.Index < src.Length {
		return src.Value[src.Index]
	}

	return NULL_CHAR
}

func (src *Source) PeekNextChar() rune {
	if src.Index+1 < src.Length {
		return src.Value[src.Index+1]
	}

	return NULL_CHAR
}
