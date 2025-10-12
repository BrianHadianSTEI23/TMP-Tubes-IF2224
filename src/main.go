/*
1. receive input from dfa reference input file and wanted pascal-s code to be compiled
2. read all the lines from the pascal-s code (read all by reading it as a turing machine / character machine)
3. for each char, the code will check and for every found certain keyword, it will be added into memory. (later, will be configured
again for the method of storage meanwhile - it's either bytecode format or filled into file once it hits a certain threshold in memory)
4. this is done until the last char of the file
*/

package main

import (
	"bufio"
	"compiler/milestone1"
	"fmt"
	"os"
	"strings"
)

func main() {

	// error handling
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panicked:", r)
		}
	}()

	// input the dfa reference file
	if len(os.Args) > 1 {
		// read the dfa reference file
		dfa_file := os.Args[1]

		// make it into data structure (using dictionary)
		dfaReference, _ := os.Open(dfa_file)
		dfaScanner := bufio.NewScanner(dfaReference)
		dfa := new(milestone1.DFA)

		for dfaScanner.Scan() {
			line := dfaScanner.Text()
			if strings.Contains(line, "Start_state") {

				// get the idx
				idx := strings.Index(line, "Start_state")

				// get the value of the start state
				dfa.StartState = line[idx:]
			} else if strings.Contains(line, "Final_state") {
				// get the idx
				idx := strings.Index(line, "Start_state")

				// get the value of the final state
				finalStates := strings.Split(line[idx:], ", ")

				dfa.FinalState = finalStates
			} else if strings.Contains(line, "State") { // assume that every other word is a state function
				elements := strings.Split(line, " ")

				transitionVal := milestone1.TransitionKey{
					State: elements[0],
					Input: elements[1],
				}

				dfa.Transition[transitionVal] = elements[2]
			}
		}

		// init variables
		currentState := dfa.StartState
		srcFile := os.Args[2]
		srcReference, _ := os.Open(srcFile)
		tokenReference, _ := os.Create("../test/output/tokens.txt")

		srcScanner := bufio.NewScanner(srcReference)
		tokenWriter := bufio.NewWriter(tokenReference)
		defer tokenWriter.Flush() // this is to make sure every transition is written into the file

		// get each text and token
		for srcScanner.Scan() {
			line := srcScanner.Text()

			// do lexical analyzer
			milestone1.LexicalAnalyzer(line, *dfa, &currentState, tokenWriter)
		}

		// final message
		fmt.Println("Tokenizing is done....")

	} else {
		fmt.Printf("Jangan lupa file DFAnya ya...")
		return
	}
}
