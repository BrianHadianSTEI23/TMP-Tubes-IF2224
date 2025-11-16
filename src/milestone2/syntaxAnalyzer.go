/**
algorithm
1. build a tree with the first node as <program>
2. while the input is not end (basically the input has not break out the recursive), do recursive on the tree
3. if some part is not right, then return special value (i do recommend epsilon) and then abort and search another branch
4. do until there is syntax error (this will return 1), and if not, it will terminate and tree and lexResult will be completed

NOTES
1. if the parser find a syntax error, it will terminate and return 1
2. if the parser find a mismatch (but there is another that matches), it will return 2
3. if the parser find that after searching through all matches and doesn't found anything, it will return 1
4. if success, it will return 0
*/

/**
TO DO
1. add production rule for type-definition
2. add production rule for parameter-group
*/

package milestone2

var productionRule = map[string][]string{
	"<program>":                 {"<program-header>", "<declaration-part>", "<compound-statement>", "DOT(.)"},
	"<program-header>":          {"KEYWORD(program)", "IDENTIFIER", "SEMICOLON(;)"},
	"<declaration-part>":        {"(const-declaration)*", "(type-declaration)*", "(var-declaration)*", "(subprogram-declaration)*"},
	"<const-declaration>":       {"KEYWORD(konstanta)", "IDENTIFIER", "=", "NUMBER", "SEMICOLON(;)"},
	"<type-declaration>":        {"KEYWORD(tipe)", "IDENTIFIER", "=", "<type-definition>", "SEMICOLON(;)"},
	"<var-declaration>":         {"KEYWORD(variabel)", "<identifier-list>", "COLON(:)", "<type>", "SEMICOLON(;)"},
	"<identifier-list>":         {"IDENTIFIER", "(COMMA(,) IDENTIFIER)*"},
	"<type>":                    {"KEYWORD(integer)", "IDENTIFIER", "SEMICOLON(;)"},
	"<array-type>":              {"KEYWORD(larik)", "LBRACKET([)", "<range>", "RBRACKET(])", "KEYWORD(dari)", "<type>"},
	"<range>":                   {"<expression>", "RANGE_OPERATOR(..)", "<expression>"},
	"<subprogram-declaration>":  {"<procedure-declaration>", "<function-declaration>"},
	"<procedure-declaration>":   {"KEYWORD(prosedur)", "IDENTIFIER", "(formal-parameter-list)*", "SEMICOLON(;)"},
	"<function-declaration>":    {"KEYWORD(function)", "IDENTIFIER", "(formal-parameter-list)*", "SEMICOLON(;)"},
	"<formal-parameter-list>":   {"LPARENTHESES(()", "<parameter-group>", "(SEMICOLON(;) <parameter-group>)*", "RPARENTHESES())"},
	"<compound-statement>":      {"KEYWORD(mulai)", "<statement-list>", "KEYWORD(selesai)"},
	"<statement>":               {"<assignment-statement>*", "<if-statement>*", "<while-statement>*", "<for-statement>*"},
	"<statement-list>":          {"<statement>", "(SEMICOLON(;) <statement>)*"},
	"<assignment-statement>":    {"IDENTIFIER", "ASSIGN-OPERATOR(:=)", "<expression>"},
	"<if-statement>":            {"KEYWORD(jika)", "<expression>", "KEYWORD(maka)", "<statement>", "(KEYWORD(selain-itu) <statement>)*"},
	"<while-statement>":         {"KEYWORD(selama)", "<expression>", "KEYWORD(lakukan)", "<statement>"},
	"<for-statement>":           {"KEYWORD(untuk)", "IDENTIFIER", "ASSIGN_OPERATOR(:=)", "<expression>", "(KEYWORD(ke) | KEYWORD(turun-ke))", "<expression>", "KEYWORD(lakukan)", "<statement>"},
	"<parameter-list>":          {"<expression>", "(COMMA(,) <expression)*"},
	"<expression>":              {"<simple-expression>", "(<relational-operator> <simple-expression>)*"},
	"<simple-expression>":       {"(ARITHMETIC_OPERATOR(+) | ARITHMETIC_OPERATOR(-))*", "<term>", "(<additive-operator> <term>)*"},
	"<term>":                    {"<factor>", "(<multiplicative-operator> <factor>)*"},
	"<factor>":                  {"(IDENTIFIER | NUMBER | CHAR_LITERAL | STRING_LITERAL | ( LPARENTHESES(() <expression> RPARENTHESES()) ) | LOGICAL_OPERATOR(tidak))", "(<factor> | <function-declaration>)"},
	"<relational-operator>":     {"= | > | < | >= | <= | <>"},
	"<additive-operator>":       {"+ | - | atau"},
	"<multiplicative-operator>": {"* | / | bagi | mod | dan"},
}

func SyntaxAnalyzer(lexResult []string, currentNode *AbstractSyntaxTree) {

	// generate many production rule that can be generated using |

	// check each process, does the element exist in the input or not

	// add the children nodes using the production rules

	// do recursive on the children nodes

	// here, the parentmost node has already finished and can be printed into txt

}
