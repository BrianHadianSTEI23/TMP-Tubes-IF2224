package milestone1

import "strings"

func Tokenize(token string) string {
	if token == "" {
		return ""
	}
	if len(token) > 0 && token[0] == '\'' && (len(token) == 1 || token[len(token)-1] != '\'') {
		return "ERROR(" + token + ")"
	}
	keywords := map[string]bool{
		"program":    true,
		"variabel":   true,
		"mulai":      true,
		"selesai":    true,
		"jika":       true,
		"maka":       true,
		"selain-itu": true,
		"selama":     true,
		"lakukan":    true,
		"untuk":      true,
		"ke":         true,
		"turun-ke":   true,
		"integer":    true,
		// "real":     true,
		"boolean":   true,
		"char":      true,
		"larik":     true,
		"dari":      true,
		"prosedur":  true,
		"fungsi":    true,
		"konstanta": true,
		"tipe":      true,
		"true":      true,
		"false":     true,
		"ulangi":    true,
		"sampai":    true,
		"kasus":     true,
		"rekaman":   true,
		"writeln":   true,
		// "read":      true,
		// "call":      true,
	}
	arithmetic := map[string]bool{
		"bagi": true,
		"mod":  true,
	}
	logical := map[string]bool{
		"dan":   true,
		"atau":  true,
		"tidak": true,
	}

	if keywords[strings.ToLower(token)] {
		return "KEYWORD(" + token + ")"
	}
	if arithmetic[strings.ToLower(token)] {
		return "ARITHMETIC_OPERATOR(" + token + ")"
	}
	if logical[strings.ToLower(token)] {
		return "LOGICAL_OPERATOR(" + token + ")"
	}
	if isNumber(token) {
		return "NUMBER(" + token + ")"
	}
	if len(token) >= 2 && token[0] == '\'' && token[len(token)-1] == '\'' {
		content := token[1 : len(token)-1]

		// char kosong
		if len(content) == 0 {
			return "CHAR_LITERAL(" + token + ")"
		}

		// len 1 = char literal
		if len(content) == 1 {
			return "CHAR_LITERAL(" + token + ")"
		}

		// else masuk ke string
		return "STRING_LITERAL(" + token + ")"
	}

	// cek operator lain berdasarkan token
	switch token {
	case "+", "-", "*", "/":
		return "ARITHMETIC_OPERATOR(" + token + ")"
	case "=", "<>", "<", ">", "<=", ">=":
		return "RELATIONAL_OPERATOR(" + token + ")"
	case ":=":
		return "ASSIGN_OPERATOR(" + token + ")"
	case ";":
		return "SEMICOLON(" + token + ")"
	case ",":
		return "COMMA(" + token + ")"
	case ":":
		return "COLON(" + token + ")"
	case ".":
		return "DOT(" + token + ")"
	case "(":
		return "LPARENTHESIS(" + token + ")"
	case ")":
		return "RPARENTHESIS(" + token + ")"
	case "[":
		return "LBRACKET(" + token + ")"
	case "]":
		return "RBRACKET(" + token + ")"
	case "..":
		return "RANGE_OPERATOR(" + token + ")"
	}

	return "IDENTIFIER(" + token + ")"
}

func isNumber(token string) bool {
	if len(token) == 0 {
		return false
	}

	for _, char := range token {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
