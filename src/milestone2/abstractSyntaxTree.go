package milestone2

import "fmt"

type AbstractSyntaxTree struct {
	value          string
	productionRule []string
	children       []*AbstractSyntaxTree
}

func printAbstractSyntaxTree(tree AbstractSyntaxTree) {
	fmt.Print(tree.value)
	if len(tree.children) != 0 {
		for i := 0; i < len(tree.children); i++ {
			printAbstractSyntaxTree((*tree.children[i]))
		}
	}
}
