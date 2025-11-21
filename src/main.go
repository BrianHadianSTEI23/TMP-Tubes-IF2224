package main

import (
	"bufio"
	"compiler/milestone1"
	"compiler/milestone2"
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

	// Cek argumen input (harus ada file DFA dan Source Code)
	if len(os.Args) < 3 {
		fmt.Printf("Cara pakai: go run ./src <file_dfa.txt> <file_program.txt>\n")
		return
	}

	dfa_file := os.Args[1]
	srcFile := os.Args[2]

	// --- 1. LOAD DFA ---
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

	// --- 2. LEXICAL ANALYZER (M1) ---
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
		// Panggil Lexer M1 (menggunakan parameter standar M1 kamu)
		milestone1.LexicalAnalyzer(line, *dfa, &currentState, tokenWriter)
	}

	// PENTING: Flush dan Close agar data tertulis sebelum dibaca M2
	tokenWriter.Flush()
	tokenReference.Close()

	// --- 3. SYNTAX ANALYZER (M2) ---
	// Inisialisasi root node untuk tree
	var root = milestone2.AbstractSyntaxTree{
		Value:    "<program>",
		Children: []*milestone2.AbstractSyntaxTree{},
	}

	// Tentukan path file tokens.txt secara absolut (biar aman)
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

	// Jalankan Syntax Analyzer (Driver M2)
	fmt.Println("Menjalankan Syntax Analysis...")
	result := milestone2.SyntaxAnalyzer(cleanLexResult, &root)

	if result == 0 { // 0 = Sukses
		fmt.Println("Syntax Analysis Berhasil! Tree dicetak ke file & terminal.")

		// Print ke Terminal (format Tree Linux)
		milestone2.PrintAbstractSyntaxTree(&root, os.Stdout, "", true)

		// Print ke File
		milestone2.PrintAbstractSyntaxTree(&root, treeWriter, "", true)
	} else {
		fmt.Println("Syntax Analysis Gagal.")
		treeWriter.WriteString("Syntax error found.")
	}

	treeWriter.Flush()
}
