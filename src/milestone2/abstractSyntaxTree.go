package milestone2

import "fmt"

type AbstractSyntaxTree struct {
	TerminalValue    string
	NonTerminalValue string
	Children         []*AbstractSyntaxTree
}

func NewNode(kind string) *AbstractSyntaxTree {
	return &AbstractSyntaxTree{NonTerminalValue: kind, Children: []*AbstractSyntaxTree{}}
}

func NewLeaf(kind, token string) *AbstractSyntaxTree {
	return &AbstractSyntaxTree{NonTerminalValue: kind, TerminalValue: token}
}

func PrintAbstractSyntaxTree(tree AbstractSyntaxTree) {
	if len(tree.TerminalValue) == 0 { // there are no terminal value (the node is a non terminal node)
		fmt.Print(tree.NonTerminalValue)
	} else {
		fmt.Print(tree.TerminalValue)
	}
	if len(tree.Children) != 0 {
		for i := 0; i < len(tree.Children); i++ {
			PrintAbstractSyntaxTree((*tree.Children[i]))
		}
	}
}
