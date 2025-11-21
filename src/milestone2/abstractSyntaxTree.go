package milestone2

import (
	"fmt"
	"io"
)

// Struktur Node Tree
type AbstractSyntaxTree struct {
	Value          string
	ProductionRule []string // (Tidak digunakan di Recursive Descent, tapi dibiarkan agar kompatibel)
	Children       []*AbstractSyntaxTree
}

// Fungsi Print Tree dengan format cantik (seperti command 'tree' di Linux)
func PrintAbstractSyntaxTree(node *AbstractSyntaxTree, writer io.Writer, prefix string, isLast bool) {
	// Cetak prefix cabang
	fmt.Fprint(writer, prefix)

	// Tentukan simbol cabang (akhir atau tengah)
	if isLast {
		fmt.Fprint(writer, "└── ")
		prefix += "    "
	} else {
		fmt.Fprint(writer, "├── ")
		prefix += "│   "
	}

	// Cetak nilai node
	fmt.Fprintln(writer, node.Value)

	// Rekursif ke anak-anak
	for i, child := range node.Children {
		PrintAbstractSyntaxTree(child, writer, prefix, i == len(node.Children)-1)
	}
}
