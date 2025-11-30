package main

import (
	"bufio"
	"compiler/milestone1"
	"compiler/milestone2"
	"compiler/milestone3"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	// Error handling agar tidak crash kotor
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panicked:", r)
		}
	}()

	// Cek argumen input
	if len(os.Args) < 3 {
		fmt.Printf("Cara pakai: go run ./src <file_dfa.txt> <file_program.txt>\n")
		return
	}

	dfa_file := os.Args[1]
	srcFile := os.Args[2]

	// 1. LOAD DFA
	dfaReference, err := os.Open(dfa_file)
	if err != nil {
		fmt.Printf("ERROR: error opening DFA file: %v\n", err)
		return
	}
	defer dfaReference.Close()

	dfaScanner := bufio.NewScanner(dfaReference)
	dfa := &milestone1.DFA{
		Transition: make(map[milestone1.TransitionKey]string),
	}

	for dfaScanner.Scan() {
		line := strings.TrimSpace(dfaScanner.Text())
		// Skip komentar dan baris kosong
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		if strings.Contains(line, "Start_state") {
			dfa.StartState = strings.TrimSpace(strings.TrimPrefix(line, "Start_state = "))
		} else if strings.Contains(line, "Final_state") {
			finalStatesStr := strings.TrimSpace(strings.TrimPrefix(line, "Final_state = "))
			finalStates := strings.Split(finalStatesStr, ", ")
			for i := range finalStates {
				finalStates[i] = strings.TrimSpace(finalStates[i])
			}
			dfa.FinalState = finalStates
		} else {
			// Baca transisi state
			elements := strings.Fields(line)
			if len(elements) >= 3 {
				transitionVal := milestone1.TransitionKey{
					State: elements[0],
					Input: elements[1],
				}
				dfa.Transition[transitionVal] = elements[2]
			}
		}
	}

	// 2. LEXICAL ANALYZER
	currentState := dfa.StartState
	srcReference, err := os.Open(srcFile)
	if err != nil {
		fmt.Printf("ERROR: error opening source file: %v\n", err)
		return
	}
	defer srcReference.Close()

	// Pastikan folder output ada
	os.MkdirAll("../test/output", os.ModePerm)
	tokenReference, err := os.Create("../test/output/tokens.txt")
	if err != nil {
		fmt.Printf("ERROR: error creating token file: %v\n", err)
		return
	}

	srcScanner := bufio.NewScanner(srcReference)
	tokenWriter := bufio.NewWriter(tokenReference)

	// Jalankan Lexer baris per baris
	for srcScanner.Scan() {
		line := srcScanner.Text()
		// Panggil Lexer
		milestone1.LexicalAnalyzer(line, *dfa, &currentState, tokenWriter)
	}

	tokenWriter.Flush()
	tokenReference.Close()

	// 3. SYNTAX ANALYZER
	// Inisialisasi root node untuk tree
	var root = milestone2.AbstractSyntaxTree{
		Value:    "<program>",
		Children: []*milestone2.AbstractSyntaxTree{},
	}

	// Tentukan path file tokens.txt secara absolut
	_, filename, _, _ := runtime.Caller(0)
	base := filepath.Dir(filename)
	path := filepath.Join(base, "..", "test", "output", "tokens.txt")

	lexResultReferenceBytes, err := os.ReadFile(path)
	if err != nil {
		// Fallback path manual jika runtime caller gagal
		path = "../test/output/tokens.txt"
		lexResultReferenceBytes, err = os.ReadFile(path)
		if err != nil {
			fmt.Printf("ERROR: Gagal membaca file tokens.txt: %v\n", err)
			return
		}
	}

	// Convert isi file jadi array string token
	lexResultStr := string(lexResultReferenceBytes)
	lexResultStr = strings.ReplaceAll(lexResultStr, "\r", "") // Bersihkan windows newline
	lexResult := strings.Split(lexResultStr, "\n")

	// Bersihkan string kosong dari array
	var cleanLexResult []string
	for _, s := range lexResult {
		if strings.TrimSpace(s) != "" {
			cleanLexResult = append(cleanLexResult, s)
		}
	}

	// Siapkan file output untuk Tree
	treeReference, err := os.Create("../test/output/abstract-syntax-tree.txt")
	if err != nil {
		fmt.Printf("ERROR: error creating tree file: %v\n", err)
		return
	}
	defer treeReference.Close()
	treeWriter := bufio.NewWriter(treeReference)

	// Jalankan Syntax Analyzer
	fmt.Println("Menjalankan Syntax Analysis...")
	result := milestone2.SyntaxAnalyzer(cleanLexResult, &root)

	if result == 0 { // 0 = Sukses
		fmt.Println("Syntax Analysis Berhasil! Tree dicetak ke file & terminal.")

		// Print ke Terminal
		milestone2.PrintAbstractSyntaxTree(&root, os.Stdout, "", true)

		// Print ke File
		milestone2.PrintAbstractSyntaxTree(&root, treeWriter, "", true)

		fmt.Println("\n========== MILESTONE 3: SEMANTIC ANALYSIS ==========")

		// Perform semantic analysis
		fmt.Println("Performing semantic analysis...")
		analyzer := milestone3.NewSemanticAnalyzer()
		decoratedAST, err := analyzer.Analyze(&root)

		if err != nil {
			fmt.Printf("Semantic analysis failed: %v\n", err)

			// Print errors
			if len(analyzer.GetErrors()) > 0 {
				fmt.Println("\nSemantic errors:")
				for i, errMsg := range analyzer.GetErrors() {
					fmt.Printf("  %d. %s\n", i+1, errMsg)
				}
			}

			// Print warnings
			if len(analyzer.GetWarnings()) > 0 {
				fmt.Println("\nSemantic warnings:")
				for i, warnMsg := range analyzer.GetWarnings() {
					fmt.Printf("  %d. %s\n", i+1, warnMsg)
				}
			}
		} else {
			fmt.Println("Semantic analysis completed successfully")

			// Print warnings if any
			if len(analyzer.GetWarnings()) > 0 {
				fmt.Println("\nSemantic warnings:")
				for i, warnMsg := range analyzer.GetWarnings() {
					fmt.Printf("  %d. %s\n", i+1, warnMsg)
				}
			}
		}

		// Get symbol table
		symTable := analyzer.GetSymbolTable()

		// Print symbol table
		fmt.Println("\n========== SYMBOL TABLE ==========")
		symTable.PrintSymbolTable()

		// Print decorated AST
		if decoratedAST != nil {
			fmt.Println("\n========== DECORATED AST ==========")
			milestone3.PrintDecoratedAST(decoratedAST, "", true)
		}

		// Save symbol table to file
		symTableFile, err := os.Create("../test/output/symbol-table.txt")
		if err == nil {
			defer symTableFile.Close()
			symTableWriter := bufio.NewWriter(symTableFile)

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			symTable.PrintSymbolTable()

			w.Close()
			os.Stdout = oldStdout

			buf := make([]byte, 1024)
			for {
				n, err := r.Read(buf)
				if n > 0 {
					symTableWriter.Write(buf[:n])
				}
				if err != nil {
					break
				}
			}
			symTableWriter.Flush()
			fmt.Println("Symbol table saved to ../test/output/symbol-table.txt")
		}

		// Save decorated AST to file
		if decoratedAST != nil {
			decoratedASTFile, err := os.Create("../test/output/decorated-ast.txt")
			if err == nil {
				defer decoratedASTFile.Close()
				decoratedASTWriter := bufio.NewWriter(decoratedASTFile)

				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				milestone3.PrintDecoratedAST(decoratedAST, "", true)

				w.Close()
				os.Stdout = oldStdout

				buf := make([]byte, 1024)
				for {
					n, err := r.Read(buf)
					if n > 0 {
						decoratedASTWriter.Write(buf[:n])
					}
					if err != nil {
						break
					}
				}
				decoratedASTWriter.Flush()
				fmt.Println("Decorated AST saved to ../test/output/decorated-ast.txt")
			}
		}
	} else {
		fmt.Println("Syntax Analysis Gagal.")
		treeWriter.WriteString("Syntax error found.")
	}

	treeWriter.Flush()
}
