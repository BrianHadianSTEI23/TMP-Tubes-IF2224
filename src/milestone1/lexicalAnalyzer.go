/*
1. for each word, it will be read from start to finish
2. if get a certain keyword, it will marked it as found and if found, next character
will be checked. if not space, it is not valid and will run into dfa until it hits a space,
and it will be assigned as identifier
*/

package milestone1

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

func LexicalAnalyzer(line string, dfa DFA, currentState *string, tokenWriter *bufio.Writer) {

	/*
		1. check if current state is already at finish
		2. if yes, then check the validity
		3. if not, go into next state based on the current input (char) and current state*/
	line = removeComments(line)

	i := 0
	for i < len(line) {
		// skip whitespaces
		for i < len(line) && unicode.IsSpace(rune(line[i])) {
			i++
		}
		if i >= len(line) {
			break
		}
		token, newPos := processToken(line, i, dfa, currentState)
		i = newPos

		if token != "" {
			convertedToken := Tokenize(token)
			if convertedToken != "" {
				tokenWriter.WriteString(convertedToken + "\n")
				fmt.Println(convertedToken)
			}
		}

	}
}

func processToken(line string, start int, dfa DFA, curr *string) (string, int) {
	*curr = dfa.StartState
	currentPos := start
	var currentToken string = ""
	longestValidToken := ""
	longestValidPos := start
	wasInString := false

	// for each char in line started from start
	for currentPos < len(line) {
		char := line[currentPos]
		if isReadingString((*curr)) {
			wasInString = true
		}
		// stop at whitespace unless string
		if unicode.IsSpace(rune(char)) && !isReadingString(*curr) {
			break
		}
		currentToken += string(char)
		// map special character
		input := mapCharForDFA(byte(char))

		transitionKey := TransitionKey{
			State: *curr,
			Input: input,
		}
		tmp, exists := dfa.Transition[transitionKey]
		if !exists {
			if wasInString {
				return currentToken, currentPos
			}
			break
		}
		*curr = tmp
		currentPos++

		// check if curr is final state
		for _, fs := range dfa.FinalState {
			if (*curr) == fs {
				longestValidToken = currentToken
				longestValidPos = currentPos
				break
			}
		}
	}
	// kasus string g ada single quote tutup
	if wasInString && longestValidToken == "" {
		return currentToken, currentPos
	}
	if longestValidToken != "" {
		return longestValidToken, longestValidPos
	}
	if start < len(line) && !unicode.IsSpace(rune(line[start])) {
		return "", start + 1
	}
	return "", start
}

// func check lg read string g
func isReadingString(state string) bool {
	return state == "STRING_START" || state == "STRING_CONTENT"
}

func removeComments(line string) string {
	// handle block comments
	for {
		start := strings.Index(line, "{")
		if start == -1 {
			break
		}
		end := strings.Index(line[start:], "}")
		if end == -1 {
			line = line[:start]
			break
		}
		line = line[:start] + line[start+end+1:]
	}
	// handle normal comments
	for {
		start := strings.Index(line, "(*")
		if start == -1 {
			break
		}
		end := strings.Index(line[start:], "*)")
		if end == -1 {
			line = line[:start]
			break
		}
		line = line[:start] + line[start+end+2:]
	}
	return line
}

func mapCharForDFA(char byte) string {
	switch char {
	case ' ':
		return "SPACE"
	case '\t':
		return "TAB"
	case '\n':
		return "NEWLINE"
	default:
		return string(char)
	}
}
