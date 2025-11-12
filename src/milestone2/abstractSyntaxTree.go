package milestone2

import "fmt"

type AbstractSyntaxTree struct {
	Value          string
	ProductionRule []string
	Children       []*AbstractSyntaxTree
}

func PrintAbstractSyntaxTree(tree AbstractSyntaxTree) {
	fmt.Print(tree.Value)
	if len(tree.Children) != 0 {
		for i := 0; i < len(tree.Children); i++ {
			PrintAbstractSyntaxTree((*tree.Children[i]))
		}
	}
}
