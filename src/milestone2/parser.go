package milestone2

import (
	"fmt"
	"regexp"
)

// Regex untuk M1 (TANPA line number)
var tokenRegexSimple = regexp.MustCompile(`^([A-Z_]+)\((.*)\)$`)

// (Struct Token ini cuma dipakai internal di M2)
type Token struct {
	Type  string
	Value string
	Line  int // (Akan 0, karena M1 apa adanya)
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%s)", t.Type, t.Value)
}

// (Struct Parser ini yang dipanggil di main.go)
type Parser struct {
	tokens  []Token
	current int
}

// (Fungsi NewParser ini yang dipanggil di main.go)
func NewParser(tokenStrings []string) *Parser {
	var tokens []Token
	for _, s := range tokenStrings {
		if s == "" {
			continue
		}

		token, err := parseTokenString(s)
		if err != nil {
			fmt.Printf("Warning: Token M1 tidak dikenali/skip: %s\n", s)
			continue
		}
		tokens = append(tokens, token)
	}
	// Tambah EOF sebagai penanda akhir
	tokens = append(tokens, Token{Type: "EOF", Value: "EOF", Line: 0})

	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

// Logic parsing string token dari file tokens.txt
func parseTokenString(s string) (Token, error) {
	matches := tokenRegexSimple.FindStringSubmatch(s)
	if matches != nil && len(matches) >= 3 {
		return Token{Type: matches[1], Value: matches[2], Line: 0}, nil
	}
	return Token{}, fmt.Errorf("invalid token format")
}

// --- Helper Functions ---

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == "EOF"
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.tokens[p.current-1]
}

func (p *Parser) check(tType string, value string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tType && p.peek().Value == value
}

func (p *Parser) checkType(tType string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tType
}

func (p *Parser) consume(tType string, value string, msg string) (*AbstractSyntaxTree, error) {
	if p.check(tType, value) {
		t := p.advance()
		return &AbstractSyntaxTree{Value: t.String()}, nil
	}
	return nil, fmt.Errorf("Syntax Error line %d: %s (Expected: %s, Got: %s(%s))", p.peek().Line, msg, value, p.peek().Type, p.peek().Value)
}

func (p *Parser) consumeType(tType string, msg string) (*AbstractSyntaxTree, error) {
	if p.checkType(tType) {
		t := p.advance()
		return &AbstractSyntaxTree{Value: t.String()}, nil
	}
	return nil, fmt.Errorf("Syntax Error line %d: %s (Got: %s)", p.peek().Line, msg, p.peek().Type)
}

// --- Recursive Descent Rules (Sesuai Grammar Spek) ---

// (Ini fungsi yang dipanggil di main.go)
func (p *Parser) ParseProgram() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<program>"}

	header, err := p.parseProgramHeader()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, header)

	decl, err := p.parseDeclarationPart()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, decl)

	compound, err := p.parseCompoundStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, compound)

	dot, err := p.consume("DOT", ".", "Expected '.' at end of program")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, dot)

	return node, nil
}

// <program-header> -> program IDENTIFIER ;
func (p *Parser) parseProgramHeader() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<program-header>"}

	prog, err := p.consume("KEYWORD", "program", "Expected 'program'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, prog)

	id, err := p.consumeType("IDENTIFIER", "Expected program name")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, id)

	semi, err := p.consume("SEMICOLON", ";", "Expected ';'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, semi)

	return node, nil
}

// <declaration-part> -> (const-decl)* (type-decl)* (var-decl)* (subprogram-decl)*
func (p *Parser) parseDeclarationPart() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<declaration-part>"}

	// (Loop untuk 'konstanta')
	for p.check("KEYWORD", "konstanta") {
		constDecl, err := p.parseConstDeclaration()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, constDecl)
	}
	// (Loop untuk 'tipe')
	for p.check("KEYWORD", "tipe") {
		// (panggil parseTypeDeclaration() nanti)
	}
	// (Loop untuk 'variabel')
	for p.check("KEYWORD", "variabel") {
		varDecl, err := p.parseVarDeclaration()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, varDecl)
	}
	// (Loop untuk 'prosedur'/'fungsi')
	for p.check("KEYWORD", "prosedur") || p.check("KEYWORD", "fungsi") {
		// (panggil parseSubprogramDeclaration() nanti)
	}
	return node, nil
}

// <const-declaration> -> konstanta (ID = value ;)+
func (p *Parser) parseConstDeclaration() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<const-declaration>"}

	kw, _ := p.consume("KEYWORD", "konstanta", "Expected 'konstanta'")
	node.Children = append(node.Children, kw)

	// (Loop untuk (...)+)
	for p.checkType("IDENTIFIER") {
		id, err := p.consumeType("IDENTIFIER", "Expected constant name")
		if err != nil {
			return nil, err
		}

		// (Spek minta '='. M1 kamu (tokenize.go) nge-token '=' sebagai RELATIONAL_OPERATOR)
		eq, err := p.consume("RELATIONAL_OPERATOR", "=", "Expected '='")
		if err != nil {
			return nil, err
		}

		// (Spek minta 'value', kita anggap NUMBER atau STRING)
		var val *AbstractSyntaxTree
		if p.checkType("NUMBER") {
			val, _ = p.consumeType("NUMBER", "")
		} else if p.checkType("STRING_LITERAL") {
			val, _ = p.consumeType("STRING_LITERAL", "")
		} else {
			return nil, fmt.Errorf("Expected NUMBER or STRING_LITERAL for constant value")
		}

		semi, err := p.consume("SEMICOLON", ";", "Expected ';'")
		if err != nil {
			return nil, err
		}

		// (Bikin sub-node biar rapi)
		constDef := &AbstractSyntaxTree{Value: "<const-def>"}
		constDef.Children = append(constDef.Children, id, eq, val, semi)
		node.Children = append(node.Children, constDef)
	}
	return node, nil
}

// <var-declaration> -> variabel (identifier-list : type ;)+
func (p *Parser) parseVarDeclaration() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<var-declaration>"}

	kw, _ := p.consume("KEYWORD", "variabel", "Expected 'variabel'")
	node.Children = append(node.Children, kw)

	for { // (Loop untuk ...)+
		idList, err := p.parseIdentifierList()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, idList)

		col, err := p.consume("COLON", ":", "Expected ':'")
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, col)

		typ, err := p.parseType()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, typ)

		semi, err := p.consume("SEMICOLON", ";", "Expected ';'")
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, semi)

		// (Jika token selanjutnya bukan ID, stop loop var-decl)
		if !p.checkType("IDENTIFIER") {
			break
		}
	}
	return node, nil
}

// <identifier-list> -> IDENTIFIER (, IDENTIFIER)*
func (p *Parser) parseIdentifierList() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<identifier-list>"}

	id, err := p.consumeType("IDENTIFIER", "Expected identifier")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, id)

	for p.check("COMMA", ",") {
		comma := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: comma.String()})

		id2, err := p.consumeType("IDENTIFIER", "Expected identifier after comma")
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, id2)
	}
	return node, nil
}

// <type> -> integer | boolean | real | char | <array-type>
func (p *Parser) parseType() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<type>"}

	if p.check("KEYWORD", "integer") ||
		p.check("KEYWORD", "boolean") ||
		p.check("KEYWORD", "real") ||
		p.check("KEYWORD", "char") {

		t := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: t.String()})
		return node, nil
	}

	if p.check("KEYWORD", "larik") {
		// (Panggil parseArrayType() nanti)
		return nil, fmt.Errorf("Parser 'larik'/'array' belum diimplementasi")
	}

	return nil, fmt.Errorf("Unknown type at line %d", p.peek().Line)
}

// <compound-statement> -> mulai <statement-list> selesai
func (p *Parser) parseCompoundStatement() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<compound-statement>"}

	start, err := p.consume("KEYWORD", "mulai", "Expected 'mulai'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, start)

	stmts, err := p.parseStatementList()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, stmts)

	end, err := p.consume("KEYWORD", "selesai", "Expected 'selesai'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, end)

	return node, nil
}

// <statement-list> -> statement (; statement)*
func (p *Parser) parseStatementList() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<statement-list>"}

	// (Handle jika blok 'mulai' kosong)
	if p.check("KEYWORD", "selesai") {
		return node, nil // Boleh kosong
	}

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, stmt)

	for p.check("SEMICOLON", ";") {
		semi := p.advance()

		// (Handle semicolon sebelum 'selesai')
		if p.check("KEYWORD", "selesai") {
			node.Children = append(node.Children, &AbstractSyntaxTree{Value: semi.String()})
			break
		}

		node.Children = append(node.Children, &AbstractSyntaxTree{Value: semi.String()})

		stmt2, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, stmt2)
	}
	return node, nil
}

// Router untuk statement
func (p *Parser) parseStatement() (*AbstractSyntaxTree, error) {
	// 1. Assignment (ID := ...)
	if p.checkType("IDENTIFIER") && p.tokens[p.current+1].Type == "ASSIGN_OPERATOR" {
		return p.parseAssignment()
	}

	// 2. Procedure Call (ID (...) )
	if p.checkType("IDENTIFIER") && p.tokens[p.current+1].Value == "(" {
		return p.parseProcedureCall()
	}

	// 3. If (jika)
	if p.check("KEYWORD", "jika") {
		return p.parseIf()
	}

	// 4. While (selama)
	if p.check("KEYWORD", "selama") {
		return p.parseWhile()
	}

	// 5. For (untuk)
	if p.check("KEYWORD", "untuk") {
		return p.parseForStatement()
	}

	// 6. Writeln (Keyword khusus)
	if p.check("KEYWORD", "writeln") {
		return p.parseProcedureCall()
	}

	// 7. Compound Nested (mulai..selesai)
	if p.check("KEYWORD", "mulai") {
		return p.parseCompoundStatement()
	}

	// (Jika tidak ada, mungkin empty, tapi kita return error jika tidak terduga)
	if p.check("KEYWORD", "selesai") {
		return &AbstractSyntaxTree{Value: "<empty-statement>"}, nil
	}

	return nil, fmt.Errorf("Unknown statement at line %d. Got: %s(%s)", p.peek().Line, p.peek().Type, p.peek().Value)
}

// <assignment> -> ID := expression
func (p *Parser) parseAssignment() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<assignment-statement>"}

	id, _ := p.consumeType("IDENTIFIER", "Expected ID")
	node.Children = append(node.Children, id)

	assign, _ := p.consume("ASSIGN_OPERATOR", ":=", "Expected :=")
	node.Children = append(node.Children, assign)

	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, expr)

	return node, nil
}

// <procedure-call> -> (ID | writeln) ( params )
func (p *Parser) parseProcedureCall() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<procedure-call>"}

	var name *AbstractSyntaxTree
	var err error
	if p.check("KEYWORD", "writeln") {
		name, err = p.consume("KEYWORD", "writeln", "")
	} else {
		// (Ini untuk prosedur buatan sendiri nanti)
		name, err = p.consumeType("IDENTIFIER", "Expected procedure name")
	}
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, name)

	// (Spek revisi 3: Kurung wajib)
	lp, err := p.consume("LPARENTHESIS", "(", "Expected '('")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, lp)

	if !p.check("RPARENTHESIS", ")") {
		params, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, params)
	}

	rp, err := p.consume("RPARENTHESIS", ")", "Expected ')'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, rp)

	return node, nil
}

// Helper untuk comma-separated expressions (untuk parameter list)
func (p *Parser) parseExprList() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<parameter-list>"}

	e1, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, e1)

	for p.check("COMMA", ",") {
		com := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: com.String()})
		e2, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, e2)
	}
	return node, nil
}

// <if-statement> -> jika expr maka stmt (selain_itu stmt)?
func (p *Parser) parseIf() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<if-statement>"}

	ifKw, _ := p.consume("KEYWORD", "jika", "")
	node.Children = append(node.Children, ifKw)

	// (Bukan mockup lagi, panggil parseExpression)
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, expr)

	thenKw, err := p.consume("KEYWORD", "maka", "Expected 'maka'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, thenKw)

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, stmt)

	if p.check("KEYWORD", "selain_itu") {
		elseKw := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: elseKw.String()})
		stmt2, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, stmt2)
	}

	return node, nil
}

// <while-statement> -> selama expr lakukan stmt
func (p *Parser) parseWhile() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<while-statement>"}

	wh, _ := p.consume("KEYWORD", "selama", "")
	node.Children = append(node.Children, wh)

	// (Bukan mockup lagi, panggil parseExpression)
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, expr)

	doKw, err := p.consume("KEYWORD", "lakukan", "Expected 'lakukan'")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, doKw)

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, stmt)

	return node, nil
}

// <for-statement> -> untuk ID := expr (ke|turun_ke) expr lakukan stmt
func (p *Parser) parseForStatement() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<for-statement>"}

	kw, _ := p.consume("KEYWORD", "untuk", "")
	node.Children = append(node.Children, kw)

	id, err := p.consumeType("IDENTIFIER", "Expected counter ID for 'for' loop")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, id)

	assign, err := p.consume("ASSIGN_OPERATOR", ":=", "Expected ':=' in 'for' loop")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, assign)

	startExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, startExpr)

	// (ke | turun_ke)
	if p.check("KEYWORD", "ke") || p.check("KEYWORD", "turun_ke") {
		dir := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: dir.String()})
	} else {
		return nil, fmt.Errorf("Expected 'ke' or 'turun_ke' in for loop")
	}

	endExpr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, endExpr)

	do, err := p.consume("KEYWORD", "lakukan", "Expected 'lakukan' in for loop")
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, do)

	stmt, err := p.parseStatement()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, stmt)

	return node, nil
}

// <expression> -> simple-expr (rel-op simple-expr)?
func (p *Parser) parseExpression() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<expression>"}

	left, err := p.parseSimpleExpression()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, left)

	if p.checkType("RELATIONAL_OPERATOR") {
		op := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: op.String()})
		right, err := p.parseSimpleExpression()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, right)
	}
	return node, nil
}

// <simple-expression> -> (+|-)? term (add-op term)*
func (p *Parser) parseSimpleExpression() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<simple-expression>"}

	// (Handle unary +/-)
	if p.check("ARITHMETIC_OPERATOR", "+") || p.check("ARITHMETIC_OPERATOR", "-") {
		op := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: op.String()})
	}

	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, left)

	for p.check("ARITHMETIC_OPERATOR", "+") || p.check("ARITHMETIC_OPERATOR", "-") || p.check("KEYWORD", "atau") {
		op := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: op.String()})
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, right)
	}
	return node, nil
}

// <term> -> factor (mul-op factor)*
func (p *Parser) parseTerm() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<term>"}

	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	node.Children = append(node.Children, left)

	for p.check("ARITHMETIC_OPERATOR", "*") || p.check("ARITHMETIC_OPERATOR", "/") ||
		p.check("ARITHMETIC_OPERATOR", "bagi") || p.check("ARITHMETIC_OPERATOR", "mod") ||
		p.check("LOGICAL_OPERATOR", "dan") {

		op := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: op.String()})
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, right)
	}
	return node, nil
}

// <factor> -> ID | ID(...) | NUM | ( expr ) | not factor
func (p *Parser) parseFactor() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<factor>"}

	// (NUMBER, STRING, CHAR, true, false)
	if p.checkType("NUMBER") || p.checkType("STRING_LITERAL") || p.checkType("CHAR_LITERAL") || p.check("KEYWORD", "true") || p.check("KEYWORD", "false") {
		t := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: t.String()})
		return node, nil
	}

	// (ID atau ID(...))
	if p.checkType("IDENTIFIER") {
		// Lookahead 1
		if p.tokens[p.current+1].Value == "(" {
			// Ini <function-call>
			funcCall, err := p.parseFunctionCall()
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, funcCall)
			return node, nil
		} else {
			// Ini cuma ID (variabel)
			t := p.advance()
			node.Children = append(node.Children, &AbstractSyntaxTree{Value: t.String()})
			return node, nil
		}
	}

	// ( <expression> )
	if p.check("LPARENTHESIS", "(") {
		lp := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: lp.String()})
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, expr)
		rp, err := p.consume("RPARENTHESIS", ")", "Expected ')'")
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, rp)
		return node, nil
	}

	// 'tidak' factor
	if p.check("LOGICAL_OPERATOR", "tidak") || p.check("KEYWORD", "tidak") {
		not := p.advance()
		node.Children = append(node.Children, &AbstractSyntaxTree{Value: not.String()})
		fact, err := p.parseFactor() // (Rekursif)
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, fact)
		return node, nil
	}

	return nil, fmt.Errorf("Unexpected token in factor: %s(%s)", p.peek().Type, p.peek().Value)
}

// <function-call> -> ID ( <expr-list> )
func (p *Parser) parseFunctionCall() (*AbstractSyntaxTree, error) {
	node := &AbstractSyntaxTree{Value: "<function-call>"} // (Sesuai spek 26)

	name, _ := p.consumeType("IDENTIFIER", "Expected function name")
	node.Children = append(node.Children, name)

	lp, _ := p.consume("LPARENTHESIS", "(", "Expected '('")
	node.Children = append(node.Children, lp)

	if !p.check("RPARENTHESIS", ")") {
		params, err := p.parseExprList()
		if err != nil {
			return nil, err
		}
		node.Children = append(node.Children, params)
	}

	rp, _ := p.consume("RPARENTHESIS", ")", "Expected ')'")
	node.Children = append(node.Children, rp)

	return node, nil
}
