package parser

import "github.com/danfragoso/milho/tokenizer"

type TokenList struct {
	Tokens []*tokenizer.Token
	Index  int
	Length int
}

func CreateTokenList(tokens []*tokenizer.Token) *TokenList {
	return &TokenList{Tokens: tokens, Length: len(tokens)}
}

func (t *TokenList) FirstToken() *tokenizer.Token {
	if t.Length > 0 {
		return t.Tokens[0]
	}

	return nil
}

func (t *TokenList) NextToken() *tokenizer.Token {
	t.Index++
	if t.Index < t.Length {
		return t.Tokens[t.Index]
	}

	return nil
}
