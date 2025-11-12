package milestone2

import (
	"fmt"
	"strings"
)

func parser(txt string) {
	for _, word := range strings.Fields(txt) {
		fmt.Println(word)
	}
}

func main() {
	txt := "duardfa adfa asdf"
	parser(txt)
}
