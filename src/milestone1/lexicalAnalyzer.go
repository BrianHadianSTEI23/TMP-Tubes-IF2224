/*
1. for each word, it will be read from start to finish
2. if get a certain keyword, it will marked it as found and if found, next character
will be checked. if not space, it is not valid and will run into dfa until it hits a space,
and it will be assigned as identifier
*/

package milestone1

import (
	"bufio"
)

func LexicalAnalyzer(line string, dfa DFA, currentState *string, tokenWriter *bufio.Writer) {

	/*
		1. check if current state is already at finish
		2. if yes, then check the validity
		3. if not, go into next state based on the current input (char) and current state*/

	var currentToken string = ""

	// first, check current is a final state?
	for _, fs := range dfa.FinalState {
		if (*currentState) == fs {
			var convertedToken string = Tokenize(currentToken)
			tokenWriter.WriteString(convertedToken)
			currentToken = ""
			break
		}
	}

	// run through for each char in the line
	for j := 0; j < len(line); j++ {

		// get the first char
		currentToken += string(line[j])
		transitionKey := TransitionKey{
			State: *currentState,
			Input: string(line[j]),
		}
		tmp := dfa.Transition[transitionKey]
		currentState = &tmp // now current state is at the second after the start state
		// (sorry about the naming because i'm too tired to think about names...)

		// check if the current state is among the final state
		for _, fs := range dfa.FinalState {
			if (*currentState) == fs {

				// if yes, check the next char whether it is a space or not
				if j+1 < len(line) {

					// if it is a space, write to buffer
					if line[j+1] == ' ' {
						// VALID! APPEND INTO OUTPUT FILE
						var convertedToken string = Tokenize(currentToken)
						tokenWriter.WriteString(convertedToken)
						currentToken = ""
					}
				}
				break
			}
		}
	}
}
