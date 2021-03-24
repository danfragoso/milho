package parser

import (
	"fmt"

	"github.com/danfragoso/milho/tokenizer"
)

func Parse(tokens []*tokenizer.Token) (*Node, error) {
	tokenList := CreateTokenList(tokens)
	currentToken := tokenList.FirstToken()

	var currentNode *Node

	for currentToken != nil {
		switch currentToken.Type {
		case tokenizer.OParen:
			if currentNode == nil {
				currentNode = createEmptyNode()
			} else {
				childNode := createEmptyNode()
				childNode.Parent = currentNode

				currentNode.Nodes = append(currentNode.Nodes, childNode)
				currentNode = childNode
			}

		case tokenizer.CParen:
			if currentNode == nil {
				return nil, fmt.Errorf("unexpected token '('")
			} else if currentNode.Parent != nil {
				currentNode = currentNode.Parent
			}

		case tokenizer.Symbol:
			if currentNode == nil {
				return nil, fmt.Errorf("unexpected token '%s'", currentToken.Value)
			} else {
				currentNode.Type = Function
				currentNode.Identifier = currentToken.Value
			}

		case tokenizer.Number:
			childNode := createEmptyNode()
			childNode.Parent = currentNode
			childNode.Type = Number
			childNode.Identifier = currentToken.Value

			if currentNode == nil {
				currentNode = childNode
			} else {
				currentNode.Nodes = append(currentNode.Nodes, childNode)
			}
		}

		currentToken = tokenList.NextToken()
	}

	return currentNode, nil
}

func createEmptyNode() *Node {
	return &Node{}
}
