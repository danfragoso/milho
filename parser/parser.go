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
	var quotes uint

	for currentToken != nil {
		quotes = incQuote(currentToken.Type, quotes)

		switch currentToken.Type {
		case tokenizer.OParen:
			childNode := createListNode()
			quotes = setQuote(childNode, quotes)

			if currentNode != nil {
				childNode.Parent = currentNode
				currentNode.Nodes = append(currentNode.Nodes, childNode)
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
			quotes = setQuote(childNode, quotes)

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

	return expandNodesSyntax(nodes)
}

func expandNodesSyntax(nodes []*Node) ([]*Node, error) {
	var nNodes []*Node
	for _, node := range nodes {
		n, err := expandNodeSyntax(node)
		if err != nil {
			return nil, err
		}

		nNodes = append(nNodes, n)
	}

	return nNodes, nil
}

func expandNodeSyntax(node *Node) (*Node, error) {
	if node.quotes > 0 {
		node = expandToQuote(node)
		return node, nil
	}

	var nodes []*Node
	for _, childNode := range node.Nodes {
		expandedNode, err := expandNodeSyntax(childNode)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, expandedNode)
	}

	node.Nodes = nodes
	return node, nil
}

func createQuoteNode() *Node {
	quoteList := createListNode()

	quoteIdentifier := createEmptyNode()
	quoteIdentifier.Type = Identifier
	quoteIdentifier.Identifier = "quote"
	quoteIdentifier.Parent = quoteList

	quoteList.Nodes = append(quoteList.Nodes, quoteIdentifier)
	return quoteList
}

func expandToQuote(node *Node) *Node {
	quoteList := createQuoteNode()
	quoteList.Parent = node.Parent

	node.Parent = quoteList
	quoteList.Nodes = append(quoteList.Nodes, node)

	for i := 1; i < int(node.quotes); i++ {
		oQuote := quoteList

		quoteList = createQuoteNode()
		quoteList.Parent = oQuote.Parent

		oQuote.Parent = quoteList
		quoteList.Nodes = append(quoteList.Nodes, oQuote)
	}

	return quoteList
}

func incQuote(tokenType tokenizer.TokenType, n uint) uint {
	if tokenType == tokenizer.SQuote {
		return n + 1
	}

	return n
}

func setQuote(node *Node, n uint) uint {
	node.quotes = n
	return 0
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
