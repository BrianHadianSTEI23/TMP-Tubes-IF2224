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

func SyntaxAnalyzer(lexResult string, currentNode AbstractSyntaxTree) {

	// push the queue (cfg's) into the current node
	if len(currentNode.children) == 0 {

		if currentNode.value == "<program>" {
			currentNode.productionRule = []string{"<program-header> <declaration-part> <compound-statement> DOT(.)"}
		} else if currentNode.value == "<program-header>" {
			currentNode.productionRule = []string{"KEYWORD(program) IDENTIFIER SEMICOLON(;)"}
		} else if currentNode.value == "<declaration-part>" {
			currentNode.productionRule = []string{"(const-declaration)* (type-declaration)* (var-declaration)* (subprogram-declaration)*"}
		} else if currentNode.value == "<const-declaration>" {
			currentNode.productionRule = []string{"KEYWORD(konstanta) IDENTIFIER = NUMBER SEMICOLON(;)"}
		} else if currentNode.value == "<type-declaration>" {
			currentNode.productionRule = []string{"KEYWORD(tipe) IDENTIFIER = <type-definition> SEMICOLON(;)"}
		} else if currentNode.value == "<var-declaration>" {
			currentNode.productionRule = []string{"KEYWORD(variabel) <identifier-list> COLON <type> SEMICOLON(;)"}
		} else if currentNode.value == "<identifier-list>" {
			currentNode.productionRule = []string{"IDENTIFIER (COMMA(,) IDENTIFIER)*"}
		} else if currentNode.value == "<type>" {
			currentNode.productionRule = []string{"KEYWORD(integer) IDENTIFIER SEMICOLON(;)"}
		} else if currentNode.value == "<array-type>" {
			currentNode.productionRule = []string{"KEYWORD(larik) LBRACKET([) <range> RBRACKET(]) KEYWORD(dari) <type>"}
		} else if currentNode.value == "<range>" {
			currentNode.productionRule = []string{"<expression> RANGE_OPERATOR(..) <expression>"}
		} else if currentNode.value == "<subprogram-declaration>" {
			currentNode.productionRule = []string{"<procedure-declaration>", "<function-declaration>"}
		} else if currentNode.value == "<procedure-declaration>" {
			currentNode.productionRule = []string{"KEYWORD(prosedur) IDENTIFIER (formal-parameter-list)* + SEMICOLON(;)"}
		} else if currentNode.value == "<function-declaration>" {
			currentNode.productionRule = []string{"KEYWORD(fungsi) IDENTIFIER (formal-parameter-list)* SEMICOLON(;)"}
		} else if currentNode.value == "<formal-parameter-list>" {
			currentNode.productionRule = []string{"LPARENTHESES(() <parameter-group> (SEMICOLON(;) <parameter-group>)* RPARENTHESES())"}
		} else if currentNode.value == "<compound-statement>" {
			currentNode.productionRule = []string{"KEYWORD(mulai) <statement-list> KEYWORD(selesai)"}
		} else if currentNode.value == "<statement>" {
			currentNode.productionRule = []string{"<assignment-statement>* <if-statement>* <while-statement>* <for-statement>*"}
		} else if currentNode.value == "<statement-list>" {
			currentNode.productionRule = []string{"<statement> (SEMICOLON(;) <statement>)*"}
		} else if currentNode.value == "<assignment-statement>" {
			currentNode.productionRule = []string{"IDENTIFIER ASSIGN-OPERATOR(:=) <expression>"}
		} else if currentNode.value == "<if-statement>" {
			currentNode.productionRule = []string{"KEYWORD(jika) <expression> KEYWORD(maka) <statement> (KEYWORD(selain-itu) <statement>)*"}
		} else if currentNode.value == "<while-statement>" {
			currentNode.productionRule = []string{"KEYWORD(selama) <expression> KEYWORD(lakukan) <statement>"}
		} else if currentNode.value == "<for-statement>" {
			currentNode.productionRule = []string{"KEYWORD(untuk) IDENTIFIER ASSIGN_OPERATOR(:=) <expression> KEYWORD(ke) <expression> KEYWORD(lakukan) <statement>", "KEYWORD(untuk) IDENTIFIER ASSIGN_OPERATOR(:=) <expression> KEYWORD(turun-ke) <expression> KEYWORD(lakukan) <statement>"}
		} else if currentNode.value == "<parameter-list>" {
			currentNode.productionRule = []string{"<expression> (COMMA(,) <expression)*"}
		} else if currentNode.value == "<expression>" {
			currentNode.productionRule = []string{"<simple-expression> (<relational-operator> <simple-expression>)*"}
		} else if currentNode.value == "<simple-expression>" {
			currentNode.productionRule = []string{"ARITHMETIC_OPERATOR(+)* <term> (<additive-operator> <term>)*"}
		} else if currentNode.value == "<term>" {
			currentNode.productionRule = []string{"<factor> (<multiplicative-operator> <factor>)*"}
		} else if currentNode.value == "<factor>" {
			currentNode.productionRule = []string{"IDENTIFIER <factor>", "NUMBER <factor>", "CHAR_LITERAL <factor>", "STRING_LITERAL <factor>", "LPARENTHESES(() <expression> RPARENTHESES()) <factor>", "LOGICAL_OPERATOR(tidak) <factor>",
				"IDENTIFIER <function-declaration>", "NUMBER <function-declaration>", "CHAR_LITERAL <function-declaration>", "STRING_LITERAL <function-declaration>", "LPARENTHESES(() <expression> RPARENTHESES()) <function-declaration>", "LOGICAL_OPERATOR(tidak) <function-declaration>"}
		} else if currentNode.value == "<relational-operator>" {
			currentNode.productionRule = []string{"=", ">", "<=", ">=", "<", "<>"}
		} else if currentNode.value == "<additive-operator>" {
			currentNode.productionRule = []string{"+", "-", "atau"}
		} else if currentNode.value == "<multiplicative-operator>" {
			currentNode.productionRule = []string{"*", "/", "bagi", "mod", "dan"}
		}
	}

	// parse the queue in the current node queue

	// add the children nodes using the production rules

	// do recursive on the children nodes
}
