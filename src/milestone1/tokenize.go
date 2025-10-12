package milestone1

func Tokenize(token string) string {
	if token == "program" { // etc
		return "KEYWORD(program)"
	}
	return "IDENTIFIER(" + token + ")"
}
