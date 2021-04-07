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
			if currentNode == nil {
				currentNode = createEmptyNode()
			} else {
				childNode := createEmptyNode()
				childNode.Parent = currentNode

				currentNode.Nodes = append(currentNode.Nodes, childNode)
				currentNode = childNode
			}

		case tokenizer.OBrack:
			if currentNode == nil {
				currentNode = createListNode()
				nodes = append(nodes, currentNode)
			} else {
				childNode := createListNode()
				childNode.Parent = currentNode

				currentNode.Nodes = append(currentNode.Nodes, childNode)
				currentNode = childNode
			}

		case tokenizer.CParen:
			if currentNode == nil {
				return nil, fmt.Errorf("unexpected token '('")
			}

			if currentNode.Parent != nil {
				currentNode = currentNode.Parent
			} else {
				nodes = append(nodes, currentNode)
				currentNode = nil
			}

		case tokenizer.CBrack:
			if currentNode == nil {
				return nil, fmt.Errorf("unexpected token '['")
			} else if currentNode.Parent != nil {
				currentNode = currentNode.Parent
			}

		case tokenizer.Reserved:
			if currentNode == nil {
				return nil, fmt.Errorf("unexpected token '%s'", currentToken.Value)
			} else {
				if currentNode.Type == Nil {
					if isMacro(currentToken.Value) {
						currentNode.Type = Macro
					} else {
						currentNode.Type = Boolean
					}

					currentNode.Identifier = currentToken.Value
				} else {
					childNode := createEmptyNode()
					childNode.Parent = currentNode

					if isMacro(currentToken.Value) {
						childNode.Type = Macro
					} else {
						childNode.Type = Boolean
					}

					childNode.Identifier = currentToken.Value

					if currentNode == nil {
						currentNode = childNode
					} else {
						currentNode.Nodes = append(currentNode.Nodes, childNode)
					}
				}
			}

		case tokenizer.Symbol:
			if currentNode == nil {
				childNode := createEmptyNode()
				childNode.Parent = currentNode
				childNode.Type = Identifier
				childNode.Identifier = currentToken.Value

				currentNode = childNode
				nodes = append(nodes, currentNode)
			} else {
				switch currentNode.Type {
				case Nil:
					currentNode.Type = Function
					currentNode.Identifier = currentToken.Value

				default:
					childNode := createEmptyNode()
					childNode.Parent = currentNode
					childNode.Type = Identifier
					childNode.Identifier = currentToken.Value

					if currentNode == nil {
						currentNode = childNode
					} else {
						currentNode.Nodes = append(currentNode.Nodes, childNode)
					}
				}
			}

		case tokenizer.Number:
			childNode := createEmptyNode()
			childNode.Parent = currentNode
			childNode.Type = Number
			childNode.Identifier = currentToken.Value

			if currentNode == nil {
				currentNode = childNode
				nodes = append(nodes, currentNode)
			} else {
				currentNode.Nodes = append(currentNode.Nodes, childNode)
			}
		}

		currentToken = tokenList.NextToken()
	}

	return nodes, nil
}

func createEmptyNode() *Node {
	return &Node{}
}

func createListNode() *Node {
	return &Node{Type: List}
}

func isMacro(macroCandidate string) bool {
	switch macroCandidate {
	case "defn", "def":
		return true
	}

	return false
}

func isBoolean(booleanCandidate string) bool {
	switch booleanCandidate {
	case "True", "False":
		return true
	}

	return false
}
