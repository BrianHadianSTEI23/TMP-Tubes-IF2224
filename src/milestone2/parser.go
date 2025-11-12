package milestone2

import (
	"fmt"
	"regexp"
	"strings"
)

func CFGToRegex(cfg string, mode string) string {
	// Escape literal parentheses
	cfg = strings.ReplaceAll(cfg, "(", `\(`)
	cfg = strings.ReplaceAll(cfg, ")", `\)`)

	// You can add more escapes here if needed
	// e.g. cfg = strings.ReplaceAll(cfg, ".", `\.`)

	// Optionally wrap with anchors
	switch mode {
	case "full":
		return "^" + cfg + "$"
	case "prefix":
		return "^" + cfg
	default: // "any"
		return cfg
	}
}

func checkNonTerminal(result string) bool {
	m, err := regexp.MatchString(`^<[A-Za-z_][A-Za-z0-9_]*`, result)

	if err != nil {
		fmt.Println("your regex is faulty")
		// you should log it or throw an error
		return false
	}
	if m {
		return true
	} else {
		return false
	}
}

func matchProductionRule(input string, currProductionRule []string, currProdRuleIndex int) string {
	i := currProdRuleIndex
	orRegExp := regexp.MustCompile("|")
	starRegExp := regexp.MustCompile(`^\(.*\)\*$`)
	for i < len(currProductionRule) {
		// if the production rule is a terminal, do check
		if checkNonTerminal(currProductionRule[i]) {

			pattern := fmt.Sprintf("^%s", currProductionRule[i]) // build pattern string
			m, err := regexp.MatchString(pattern, input)
			if !m { // not found any matching string in the prod rule
				return ""
			}
			if err != nil { // the regex is wrong
				fmt.Println("your regex is faulty")
				return ""
			}
			if i < currProdRuleIndex { // the found product rule is behind the currPrdRuleIndex, continue
				continue
			}

			// examine the productionRule thoroughly and compare it with the input

			// check if there is | delimiter, ()* delimiter
			orDelimiter := orRegExp.MatchString(currProductionRule[i])
			starDelimiter := starRegExp.MatchString(currProductionRule[i])
			if orDelimiter || starDelimiter {

				m, err := regexp.MatchString(CFGToRegex(currProductionRule[i], "prefix"), input)
				if err != nil { // the regex is wrong
					fmt.Println("your regex is faulty")
					return ""
				}
				if m { // found matching string in the prod rule
					return input // this is a workaround that it will return the right value of the matching string
				} else {
					return ""
				}

			}

			// if there are no delimiter etc etc, it will return the production rule if it's the same as input
			terminal, _ := regexp.MatchString(CFGToRegex(currProductionRule[i], "prefix"), input)
			if terminal { // if the input is the same like the production rule
				return input
			}
			i++ // go to next part of prod rule
		} else { // if not a terminal, then expand
			currProductionRule = append(currProductionRule[:i], append(productionRule[currProductionRule[i]], currProductionRule[i:]...)...) // basically this expand the currProductionRule (i'm sorry it's gibberish)
		}
	}
	return currProductionRule[i]
}
