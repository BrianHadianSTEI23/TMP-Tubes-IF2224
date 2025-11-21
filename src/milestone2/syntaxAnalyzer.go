package milestone2

import (
	"fmt"
	"regexp"
)

// Regex untuk membaca format token dari output M1 (KEYWORD(program))
var tokenParserRegex = regexp.MustCompile(`^([A-Z_]+)\((.*)\)$`)

// Wrapper function yang dipanggil main.go
// Ini menghubungkan data string dari main.go ke logika Parser baru
func SyntaxAnalyzer(lexResult []string, rootNode *AbstractSyntaxTree) int {

	// 1. Konversi String -> Struct Token
	var tokens []Token
	for _, s := range lexResult {
		if s == "" {
			continue
		}

		matches := tokenParserRegex.FindStringSubmatch(s)
		if matches != nil && len(matches) >= 3 {
			tokens = append(tokens, Token{Type: matches[1], Value: matches[2]})
		} else {
			// Jika format tidak dikenali (misal error message), skip atau handle
			// fmt.Printf("Warning: Token skipped: %s\n", s)
		}
	}

	// Tambah token EOF di akhir untuk menandakan selesai
	tokens = append(tokens, Token{Type: "EOF", Value: "EOF"})

	// 2. Setup Parser
	p := &Parser{
		tokens:  tokens,
		current: 0,
	}

	// 3. Mulai Parsing
	parsedNode, err := p.ParseProgram() 

	if err != nil {
		fmt.Printf("\n[Syntax Error] %v\n", err)
		return 1 // Return 1 menandakan error
	}

	*rootNode = *parsedNode

	return 0 // Return 0 = Sukses
}
