/*
1. for each word, it will be read from start to finish
2. if get a certain keyword, it will marked it as found and if found, next character
will be checked. if not space, it is not valid and will run into dfa until it hits a space,
and it will be assigned as identifier
*/

package milestone1

func LexicalAnalyzer(line string, dfa DFA, currentState *string) {

	/*
		1. check if current state is already at finish
		2. if yes, then check the validity
		3. if not, go into next state based on the current input (char) and current state*/

	if keys[iter] == string(line[j]) {
		found = true
		iter = 0
		// check if found, will the next char be space?
		if j+1 < len(line) {
			if line[j+1] == ' ' {
				// VALID! APPEND INTO OUTPUT FILE
				writer.WriteString(dfa.FinalState[keys[iter]])
			} else {
				// INVALID! CHECK AGAIN
				writer.WriteString(dfa.FinalState[keys[iter]])

			}
		}
	} else {
		iter++
	}

	// run through for each char in the line
	for j := 0; j < len(line); j++ {

		// check if the char is a final state, FOUND!
		found := false
		iter := 0

		for !found && iter < len(keys) {

			if keys[iter] == string(line[j]) {
				found = true
				iter = 0
				// check if found, will the next char be space?
				if j+1 < len(line) {
					if line[j+1] == ' ' {
						// VALID! APPEND INTO OUTPUT FILE
						writer.WriteString(dfa.FinalState[keys[iter]])
					} else {
						// INVALID! CHECK AGAIN
						writer.WriteString(dfa.FinalState[keys[iter]])

					}
				}
			} else {
				iter++
			}

		}

		// if stil not found, then check using the dfa file

	}
}
