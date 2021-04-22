package parser

import (
	"fmt"

	"github.com/danfragoso/milho/tokenizer"
)

func Parse(tokens []*tokenizer.Token) ([]*Node, error) {
	tokenList := CreateTokenList(tokens)
	currentToken := tokenList.FirstToken()

	var currentNode *Node
	var nodes []*Node

	for currentToken != nil {
		switch currentToken.Type {
		case tokenizer.OParen:
			childNode := createListNode()
			if currentNode != nil {
				childNode.Parent = currentNode
				currentNode.Nodes = append(currentNode.Nodes, childNode)
			}

			if tokenList.PreviousToken() != nil &&
				tokenList.PreviousToken().Type == tokenizer.SQuote {
				childNode.notToeval = true
			}

			currentNode = childNode

		case tokenizer.CParen:
			if currentNode == nil {
				return nil, fmt.Errorf("Unmatched delimiter ')'")
			}

			if currentNode.Parent != nil {
				currentNode = currentNode.Parent
			} else {
				nodes = append(nodes, currentNode)
				currentNode = nil
			}

		case tokenizer.Symbol, tokenizer.Number, tokenizer.Boolean, tokenizer.String:
			childNode := createEmptyNode()
			if currentNode != nil {
				childNode.Parent = currentNode
				currentNode.Nodes = append(currentNode.Nodes, childNode)
			} else {
				nodes = append(nodes, childNode)
			}

			childNode.Type = resolveNodeType(currentToken)
			childNode.Identifier = currentToken.Value
		}

		currentToken = tokenList.NextToken()
	}

	return nodes, nil
}

func resolveNodeType(token *tokenizer.Token) NodeType {
	switch token.Type {
	case tokenizer.Symbol:
		return Identifier

	case tokenizer.Number:
		return Number

	case tokenizer.Boolean:
		return Boolean

	case tokenizer.String:
		return String
	}

	return Nil
}

func createEmptyNode() *Node {
	return &Node{}
}

func createListNode() *Node {
	return &Node{
		Type: List,
	}
}

func isBoolean(booleanCandidate string) bool {
	switch booleanCandidate {
	case "True", "False":
		return true
	}

	return false
}
