// this file basically fill all up with the cfg and others

package milestone2

import (
	"fmt"
	"regexp"
)

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

type Parser struct {
	Tokens []string
	Pos    int
}

func NewParser(Tokens []string) *Parser {
	return &Parser{Tokens: Tokens, Pos: 0}
}

// Helpers
func (p *Parser) eof() bool {
	return p.Pos >= len(p.Tokens)
}
func (p *Parser) peek() string {
	if p.eof() {
		return ""
	}
	return p.Tokens[p.Pos]
}
func (p *Parser) next() string {
	if p.eof() {
		return ""
	}
	t := p.Tokens[p.Pos]
	p.Pos++
	return t
}
func (p *Parser) accept(tok string) bool {
	m, _ := regexp.MatchString(`^`+tok, p.peek())
	if m {
		p.next()
		return true
	}
	return false
}
func (p *Parser) expect(tok string) error {
	m, _ := regexp.MatchString(`^`+tok, p.peek())
	if m {
		p.next()
		return nil
	}
	return fmt.Errorf("expected %q but got %q at Pos %d", tok, p.peek(), p.Pos)
}

func (p *Parser) ParseProgram() (*AbstractSyntaxTree, error) {
	root := NewNode("<program>") // entrypoint for every program

	// append the program-header if valid
	ph, err := p.ParseProgramHeader()
	if err != nil {
		return (nil), err
	}
	root.Children = append(root.Children, ph)

	// append the declaration-part if valid
	dp, err := p.ParseDeclarationPart()
	if err != nil {
		return (nil), err
	}
	root.Children = append(root.Children, dp)

	// append the program-header if valid
	cs, err := p.ParseCompoundStatement()
	if err != nil {
		return (nil), err
	}
	root.Children = append(root.Children, cs)

	// append the program-header if valid
	derr := p.expect("^DOT(.)")
	if derr != nil {
		return (nil), derr
	}
	root.Children = append(root.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return root, nil
}

func (p *Parser) ParseProgramHeader() (*AbstractSyntaxTree, error) {
	// create program header node and if exist, append it to the children
	ph := NewNode("<program-header>")
	var err error

	// expect KEYWORD(program)
	err = p.expect("KEYWORD(program)")
	if err != nil {
		return nil, err
	}
	ph.Children = append(ph.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTIFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	ph.Children = append(ph.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect SEMICOLON(;)
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	ph.Children = append(ph.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return ph, nil
}

// "<declaration-part>":        {"(const-declaration)*", "(type-declaration)*", "(var-declaration)*", "(subprogram-declaration)*"},
func (p *Parser) ParseDeclarationPart() (*AbstractSyntaxTree, error) {
	// init + create main node
	dp := NewNode("<declaration-part>")

	// accept const declaration
	for {
		if p.peek() == "KEYWORD(konstanta)" {
			cd, err := p.ParseConstDeclaration()
			if err != nil {
				return nil, err
			}
			dp.Children = append(dp.Children, cd)
			continue
		}

		if p.peek() == "KEYWORD(variabel)" {
			vd, err := p.ParseVarDeclaration()
			if err != nil {
				return nil, err
			}
			dp.Children = append(dp.Children, vd)
			continue
		}

		if p.peek() == "KEYWORD(tipe)" {
			td, err := p.ParseTypeDeclaration()
			if err != nil {
				return nil, err
			}
			dp.Children = append(dp.Children, td)
			continue
		}

		if p.peek() == "KEYWORD(procedure)" {
			pd, err := p.ParseProcedureDeclaration()
			if err != nil {
				return nil, err
			}
			dp.Children = append(dp.Children, pd)
			continue
		}

		if p.peek() == "KEYWORD(function)" {
			fd, err := p.ParseFunctionDeclaration()
			if err != nil {
				return nil, err
			}
			dp.Children = append(dp.Children, fd)
			continue
		}
		break
	}
	return dp, nil
}

// "<const-declaration>":       {"KEYWORD(konstanta)", "IDENTIFIER", "=", "NUMBER", "SEMICOLON(;)"},
func (p *Parser) ParseConstDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	cd := NewNode("<const-declaration>")
	var err error

	// expect KEYWORD(program)
	err = p.expect("KEYWORD(konstanta)")
	if err != nil {
		return nil, err
	}
	cd.Children = append(cd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTIFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	cd.Children = append(cd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect assign operator
	err = p.expect("ASSIGN_OPERATOR(=)")
	if err != nil {
		return nil, err
	}
	cd.Children = append(cd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect number
	err = p.expect("NUMBER")
	if err != nil {
		return nil, err
	}
	cd.Children = append(cd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect semicolon
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	cd.Children = append(cd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return cd, nil
}

// "<type-declaration>":        {"KEYWORD(tipe)", "IDENTIFIER", "=", "<type-definition>", "SEMICOLON(;)"},
func (p *Parser) ParseTypeDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	td := NewNode("<type-declaration>")
	var err error

	// expect KEYWORD(tipe)
	err = p.expect("KEYWORD(tipe)")
	if err != nil {
		return nil, err
	}
	td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTIFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect assign operator
	err = p.expect("ASSIGN_OPERATOR(=)")
	if err != nil {
		return nil, err
	}
	td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect type definitiona
	for {
		n, nerr := regexp.MatchString("^NUMBER", p.peek())
		c, cerr := regexp.MatchString("^CHAR_LITERAL", p.peek())
		s, serr := regexp.MatchString("^STRING_LITERAL", p.peek())
		if n {
			_ = p.accept("NUMBER")
			if nerr != nil {
				return nil, nerr
			}
			td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable
		} else if c {
			_ = p.accept("CHAR_LITERAL")
			if cerr != nil {
				return nil, cerr
			}
			td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable
		} else if s {
			_ = p.accept("STRING_LITERAL")
			if serr != nil {
				return nil, serr
			}
			td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable
		}
		break
	}

	// expect semicolon
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return td, nil
}

// "<var-declaration>":         {"KEYWORD(variabel)", "<identifier-list>", "COLON(:)", "<type>", "SEMICOLON(;)"},
func (p *Parser) ParseVarDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	vd := NewNode("<var-declaration>")
	var err error

	// expect KEYWORD(program)
	err = p.expect("KEYWORD(variabel)")
	if err != nil {
		return nil, err
	}
	vd.Children = append(vd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <identifier-list>
	il, err := p.ParseIdentifierList()
	if err != nil {
		return (nil), err
	}
	vd.Children = append(vd.Children, il)

	// expect COLON(:)
	err = p.expect("COLON(:)")
	if err != nil {
		return nil, err
	}
	vd.Children = append(vd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	// expect type
	t, err := p.ParseType()
	if err != nil {
		return (nil), err
	}
	vd.Children = append(vd.Children, t)

	// expect semicolon
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	vd.Children = append(vd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return vd, nil
}

func (p *Parser) ParseIdentifierList() (*AbstractSyntaxTree, error) {
	// init + create main node
	il := NewNode("<identifier-list>")
	var err error

	// expect semicolon
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	il.Children = append(il.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// accept comma(,) IDENTIFIER
	for {
		if p.peek() == "COMMA(,)" {
			err = p.expect("COMMA(,)")
			if err != nil {
				return nil, err
			}
			il.Children = append(il.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

			// expect IDENTIFIER
			_, ierr := regexp.MatchString("^IDENTIFIER", p.peek())

			if ierr != nil {
				return nil, err
			}

			il.Children = append(il.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
			continue
		}

		break
	}

	return il, nil
}

func (p *Parser) ParseType() (*AbstractSyntaxTree, error) {
	// init + create main node
	t := NewNode("<type>")
	var err error

	// expect KEYWORD(integer)
	err = p.expect("KEYWORD(integer)")
	if err != nil {
		return nil, err
	}
	t.Children = append(t.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTIFIER
	err = p.expect("IDENTIFIER")
	// il, err := p.ParseIdentifierList()
	if err != nil {
		return (nil), err
	}
	t.Children = append(t.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect SEMICOLON(;)
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	t.Children = append(t.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable

	return t, nil
}

// 	"<array-type>":              {"KEYWORD(larik)", "LBRACKET([)", "<range>", "RBRACKET(])", "KEYWORD(dari)", "<type>"},

func (p *Parser) ParseArrayType() (*AbstractSyntaxTree, error) {
	// init + create main node
	at := NewNode("<array-type>")
	var err error

	// expect KEYWORD(larik)
	err = p.expect("KEYWORD(larik)")
	if err != nil {
		return nil, err
	}
	at.Children = append(at.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect LBRACKET([)
	err = p.expect("LBRACKET([)")
	if err != nil {
		return nil, err
	}
	at.Children = append(at.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <range>
	r, err := p.ParseRange()
	if err != nil {
		return (nil), err
	}
	at.Children = append(at.Children, r)

	// expect RBRACKET(])
	err = p.expect("RBRACKET(])")
	if err != nil {
		return nil, err
	}
	at.Children = append(at.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect KEYWORD(dari)
	err = p.expect("KEYWORD(dari)")
	if err != nil {
		return nil, err
	}
	at.Children = append(at.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <type>
	t, terr := p.ParseType()
	// N(:)")
	if terr != nil {
		return nil, terr
	}
	at.Children = append(at.Children, t) // this is minus one because p.expect has incremented the Pos variable

	return at, nil
}

// "<range>":                   {"<expression>", "RANGE_OPERATOR(..)", "<expression>"},
func (p *Parser) ParseRange() (*AbstractSyntaxTree, error) {
	// init + create main node
	r := NewNode("<range>")
	var err error

	// expect <expression>
	e, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	r.Children = append(r.Children, e)

	// expect RANGE_OPERATOR(..)
	err = p.expect("RANGE_OPERATOR(..)")
	if err != nil {
		return nil, err
	}
	r.Children = append(r.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <expression>
	e, eerr := p.ParseExpression()
	// N(:)")
	if eerr != nil {
		return nil, eerr
	}
	r.Children = append(r.Children, e) // this is minus one because p.expect has incremented the Pos variable

	return r, nil
}

// "<subprogram-declaration>":  {"<procedure-declaration>", "<function-declaration>"},
func (p *Parser) ParseSubprogramDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	sd := NewNode("<subprogram-declaration>")

	for {
		if p.peek() == "KEYWORD(prosedur)" {
			pd, pderr := p.ParseProcedureDeclaration()
			if pderr != nil {
				return nil, pderr
			}
			sd.Children = append(sd.Children, pd)
		} else if p.peek() == "KEYWORD(fungsi)" {
			fd, fderr := p.ParseFunctionDeclaration()
			if fderr != nil {
				return nil, fderr
			}
			sd.Children = append(sd.Children, fd)
		}
		break
	}

	return sd, nil
}

// "<function-declaration>":    {"KEYWORD(fungsi)", "IDENTIFIER", "(formal-parameter-list)*", "SEMICOLON(;)"},
func (p *Parser) ParseFunctionDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	fd := NewNode("<function-declaration>")
	var err error

	// expect KEYWORD(function)
	err = p.expect("KEYWORD(fungsi)")
	if err != nil {
		return nil, err
	}
	fd.Children = append(fd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	fd.Children = append(fd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// accept (<formal-parameter-list)*
	for {
		if p.peek() == "LPARENTHESES(()" {
			fpl, fplerr := p.ParseFormalParameterList()
			if fplerr != nil {
				return nil, fplerr
			}
			fd.Children = append(fd.Children, fpl)
			continue
		}
		break
	}

	// expect SEMICOLON(;)
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	fd.Children = append(fd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return fd, nil
}

func (p *Parser) ParseProcedureDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	pd := NewNode("<procedure-declaration>")
	var err error

	// expect KEYWORD(prosedur)
	err = p.expect("KEYWORD(prosedur)")
	if err != nil {
		return nil, err
	}
	pd.Children = append(pd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	pd.Children = append(pd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// accept (<formal-parameter-list)*
	for {
		if p.peek() == "LPARENTHESES(()" {
			fpl, fplerr := p.ParseFormalParameterList()
			if fplerr != nil {
				return nil, fplerr
			}
			pd.Children = append(pd.Children, fpl)
			continue
		}
		break
	}

	// expect SEMICOLON(;)
	err = p.expect("SEMICOLON(;)")
	if err != nil {
		return nil, err
	}
	pd.Children = append(pd.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return pd, nil
}

// <parameter-group> : IDENTIFIER (COMMA(,) IDENTIFIER)* COLON(:) <type>
func (p *Parser) ParseParameterGroup() (*AbstractSyntaxTree, error) {
	// init + create main node
	pg := NewNode("<parameter-group>")
	var err error

	// expaect IDENTIFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	pg.Children = append(pg.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// accept comma(,) IDENTIFIER
	for {
		if p.peek() == "COMMA(,)" {
			err = p.expect("COMMA(,)")
			if err != nil {
				return nil, err
			}
			pg.Children = append(pg.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

			// expect IDENTIFIER
			_, ierr := regexp.MatchString("^IDENTIFIER", p.peek())

			if ierr != nil {
				return nil, err
			}

			pg.Children = append(pg.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
			continue
		}

		break
	}

	// expect COLON(:)
	err = p.expect("COLON(:)")
	if err != nil {
		return nil, err
	}
	pg.Children = append(pg.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect SEMICOLON(;)
	t, terr := p.ParseType()
	if terr != nil {
		return nil, terr
	}
	t.Children = append(t.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return t, nil
}

// "<formal-parameter-list>":   {"LPARENTHESES(()", "<parameter-group>", "(SEMICOLON(;) <parameter-group>)*", "RPARENTHESES())"},
func (p *Parser) ParseFormalParameterList() (*AbstractSyntaxTree, error) {
	// init + create main node
	fpl := NewNode("<formal-parameter-list>")
	var err error

	// expect LPARENTHESES(()
	err = p.expect("LPARENTHESES(()")
	if err != nil {
		return nil, err
	}
	fpl.Children = append(fpl.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <parameter-group>
	pg, pgerr := p.ParseParameterGroup()
	if pgerr != nil {
		return nil, pgerr
	}
	fpl.Children = append(fpl.Children, pg)

	// accept (SEMICOLON(;) <parameter-group>)*
	for {
		if p.peek() == "SEMICOLON(;)" {
			err = p.expect("SEMICOLON(;)")
			if err != nil {
				return nil, err
			}
			fpl.Children = append(fpl.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

			// expect <parameter-group>
			pg, pgerr := p.ParseParameterGroup()

			if pgerr != nil {
				return nil, pgerr
			}

			fpl.Children = append(pg.Children, pg)
			continue
		}

		break
	}

	// expect RPARENTHESES())
	err = p.expect("RPARENTHESES())")
	if err != nil {
		return nil, err
	}
	fpl.Children = append(fpl.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return fpl, nil
}

// "<compound-statement>":      {"KEYWORD(mulai)", "<statement-list>", "KEYWORD(selesai)"},
func (p *Parser) ParseCompoundStatement() (*AbstractSyntaxTree, error) {
	// init + create main node
	cs := NewNode("<compoung-statement>")
	var err error

	// expect KEYWORD(mulai)
	err = p.expect("KEYWORD(mulai)")
	if err != nil {
		return nil, err
	}
	cs.Children = append(cs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <statement-list>
	sl, slerr := p.ParseStatementList()
	if slerr != nil {
		return nil, slerr
	}
	cs.Children = append(cs.Children, sl)

	// expect KEYWORD(SELESAI)
	err = p.expect("KEYWORD(selesai)")
	if err != nil {
		return nil, err
	}
	cs.Children = append(cs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	return cs, nil
}

// /////////////////////////////////////////////////// INI KEBAWAH BELUM DIBENERIN /////////////////////////
//
//	"<statement>":               {"<assignment-statement>*", "<if-statement>*", "<while-statement>*", "<for-statement>*"},
func (p *Parser) ParseStatement() (*AbstractSyntaxTree, error) {
	// init
	s := NewNode("<statement>")

	// accept const declaration
	for {
		_, ierr := regexp.MatchString("^IDENTIFIER", p.peek())
		if ierr == nil {
			as, aserr := p.ParseAssignmentStatement()
			if aserr != nil {
				return nil, aserr
			}
			s.Children = append(s.Children, as)
			continue
		}

		if p.peek() == "KEYWORD(jika)" {
			is, err := p.ParseIfStatement()
			if err != nil {
				return nil, err
			}
			s.Children = append(s.Children, is)
			continue
		}

		if p.peek() == "KEYWORD(selama)" {
			ws, err := p.ParseWhileStatement()
			if err != nil {
				return nil, err
			}
			s.Children = append(s.Children, ws)
			continue
		}

		if p.peek() == "KEYWORD(untuk)" {
			ws, err := p.ParseWhileStatement()
			if err != nil {
				return nil, err
			}
			s.Children = append(s.Children, ws)
			continue
		}

		break
	}

	return s, nil
}

// "<statement-list>":          {"<statement>", "(SEMICOLON(;) <statement>)*"},
func (p *Parser) ParseStatementList() (*AbstractSyntaxTree, error) {
	// init + create main node
	sl := NewNode("<statement-list>")
	var err error

	// expect <statement>
	s, serr := p.ParseStatement()
	if serr != nil {
		return nil, serr
	}
	sl.Children = append(sl.Children, s)

	// accept (SEMICOLON(;) <statement>)*
	for {
		if p.peek() == "SEMICOLON(;)" {
			err = p.expect("SEMICOLON(;)")
			if err != nil {
				return nil, err
			}
			sl.Children = append(sl.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

			// expect <staement>
			s, serr := p.ParseStatement()

			if serr != nil {
				return nil, serr
			}

			sl.Children = append(sl.Children, s)
			continue
		}

		break
	}
	return sl, nil
}

// "<assignment-statement>":    {"IDENTIFIER", "ASSIGN-OPERATOR(:=)", "<expression>"},
func (p *Parser) ParseAssignmentStatement() (*AbstractSyntaxTree, error) {
	// init + create main node
	as := NewNode("<assignment-statement>")
	var err error

	// expect IDENTIFIER
	err = p.expect("IDENTFIER")
	if err != nil {
		return nil, err
	}
	as.Children = append(as.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect ASSIGN_OPERATOR(:=)
	err = p.expect("ASSIGN_OPERATOR(:=)")
	if err != nil {
		return nil, err
	}
	as.Children = append(as.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <expression>
	e, eerr := p.ParseExpression()
	if eerr != nil {
		return (nil), eerr
	}
	as.Children = append(as.Children, e)

	return as, nil
}

// "<if-statement>":            {"KEYWORD(jika)", "<expression>", "KEYWORD(maka)", "<statement>", "(KEYWORD(selain-itu) <statement>)*"},
func (p *Parser) ParseIfStatement() (*AbstractSyntaxTree, error) {
	// init + create main node
	is := NewNode("<if-statement>")
	var err error

	// expect KEYWORD(jika)
	err = p.expect("KEYWORD(jika)")
	if err != nil {
		return nil, err
	}
	is.Children = append(is.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <expresison>
	e, eerr := p.ParseExpression()
	if eerr != nil {
		return (nil), eerr
	}
	is.Children = append(is.Children, e)

	// expect KEYWORD(maka)
	err = p.expect("KEYWORD(maka)")
	if err != nil {
		return nil, err
	}
	is.Children = append(is.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <statement>
	s, serr := p.ParseStatement()
	if serr != nil {
		return (nil), serr
	}
	is.Children = append(is.Children, s)

	// accept (SEMICOLON(;) <statement>)*
	for {
		if p.peek() == "SEMICOLON(;)" {
			err = p.expect("SEMICOLON(;)")
			if err != nil {
				return nil, err
			}
			is.Children = append(is.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

			// expect <staement>
			s, serr := p.ParseStatement()

			if serr != nil {
				return nil, serr
			}

			is.Children = append(is.Children, s)
			continue
		}

		break
	}

	return is, nil
}

// "<while-statement>":         {"KEYWORD(selama)", "<expression>", "KEYWORD(lakukan)", "<statement>"},
func (p *Parser) ParseWhileStatement() (*AbstractSyntaxTree, error) {
	// init + create main node
	ws := NewNode("<while-statement>")
	var err error

	// expect KEYWORD(selama)
	err = p.expect("KEYWORD(selama)")
	if err != nil {
		return nil, err
	}
	ws.Children = append(ws.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <expression>
	e, eerr := p.ParseExpression()
	if eerr != nil {
		return (nil), eerr
	}
	ws.Children = append(ws.Children, e)

	// expect KEYWORD(lakukan)
	err = p.expect("KEYWORD(lakukan)")
	if err != nil {
		return nil, err
	}
	ws.Children = append(ws.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <statement>
	s, serr := p.ParseStatement()
	if serr != nil {
		return (nil), serr
	}
	ws.Children = append(ws.Children, s)

	return ws, nil
}

// "<for-statement>":           {"KEYWORD(untuk)", "IDENTIFIER", "ASSIGN_OPERATOR(:=)", "<expression>", "(KEYWORD(ke) | KEYWORD(turun-ke))", "<expression>", "KEYWORD(lakukan)", "<statement>"},
func (p *Parser) ParseForStatement() (*AbstractSyntaxTree, error) {
	// init + create main node
	fs := NewNode("<for-statement>")
	var err error

	// expect KEYWORD(untuk)
	err = p.expect("KEYWORD(untuk)")
	if err != nil {
		return nil, err
	}
	fs.Children = append(fs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect IDENTIFIER
	err = p.expect("IDENTIFIER")
	if err != nil {
		return nil, err
	}
	fs.Children = append(fs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect ASSIGN_OPERATOR(:=)
	err = p.expect("ASSIGN_OPERATOR(:=)")
	if err != nil {
		return nil, err
	}
	fs.Children = append(fs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <expresison>
	e, eerr := p.ParseExpression()
	if eerr != nil {
		return (nil), eerr
	}
	fs.Children = append(fs.Children, e)

	// accept (KEYWORD(ke) | KEYWORD(turun-ke))*
	for {
		if p.peek() == "KEYWORD(ke)" {
			err = p.expect("KEYWORD(ke)")
			if err != nil {
				return nil, err
			}
			fs.Children = append(fs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		} else if p.peek() == "KEYWORD(turun-ke)" {
			err = p.expect("KEYWORD(turun-ke)")
			if err != nil {
				return nil, err
			}
			fs.Children = append(fs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		}
		break
	}

	// expect <expression>
	e, eerr = p.ParseExpression()
	if eerr != nil {
		return (nil), eerr
	}
	fs.Children = append(fs.Children, e)

	// expect KEYWORD(lakukan)
	err = p.expect("KEYWORD(lakukan)")
	if err != nil {
		return nil, err
	}
	fs.Children = append(fs.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

	// expect <statement>
	s, serr := p.ParseStatement()
	if serr != nil {
		return (nil), serr
	}
	fs.Children = append(fs.Children, s)

	return fs, nil
}

// "<parameter-list>":          {"<expression>", "(COMMA(,) <expression)*"},
func (p *Parser) ParseParameterList() (*AbstractSyntaxTree, error) {
	// init + create main node
	pl := NewNode("<parameter-list>")
	var err error

	// expect <expression>
	e, eerr := p.ParseExpression()
	if eerr != nil {
		return (nil), eerr
	}
	pl.Children = append(pl.Children, e)

	// accept (COMMA(;) <expression>)*
	for {
		if p.peek() == "COLON(:)" {
			err = p.expect("COLON(;)")
			if err != nil {
				return nil, err
			}
			pl.Children = append(pl.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))

			// expect <expression>
			s, serr := p.ParseExpression()

			if serr != nil {
				return nil, serr
			}

			pl.Children = append(pl.Children, s)
			continue
		}

		break
	}
	return pl, nil
}

// "<expression>":              {"<simple-expression>", "(<relational-operator> <simple-expression>)*"},
func (p *Parser) ParseExpression() (*AbstractSyntaxTree, error) {
	// init + create main node
	e := NewNode("<expression>")

	// expect <simple-expression>
	se, seerr := p.ParseSimpleExpression()
	if seerr != nil {
		return nil, seerr
	}
	e.Children = append(e.Children, se)

	// accept (<relational-operator> <simple-expression>)*
	for {
		ro, _ := regexp.MatchString("^RELATIONAL_OPERATOR", p.peek())
		if ro {
			ro, roerr := p.ParseRelationalOperator()
			if roerr != nil {
				return nil, roerr
			}
			e.Children = append(e.Children, ro)

			// expect <simple-expression>
			se, seerr = p.ParseSimpleExpression()

			if seerr != nil {
				return nil, seerr
			}

			e.Children = append(e.Children, se)
			continue
		}

		break
	}

	return e, nil
}

// "<simple-expression>":       {"(ARITHMETIC_OPERATOR(+) | ARITHMETIC_OPERATOR(-))*", "<term>", "(<additive-operator> <term>)*"}
func (p *Parser) ParseSimpleExpression() (*AbstractSyntaxTree, error) {
	// init + create main node
	se := NewNode("<simple-expression>")

	// accept (ARITHMETIC_OPERATOR(+) | ARITHMETIC_OPERATOR(-))*
	ro, roerr := p.ParseRelationalOperator()
	if roerr != nil {
		return nil, roerr
	}
	se.Children = append(se.Children, ro)

	// expect <term>
	t, terr := p.ParseTerm()
	if terr != nil {
		return nil, terr
	}
	se.Children = append(se.Children, t)

	// accept (<additive-operator> <term>)*
	for {
		aro, _ := regexp.MatchString("^ARITHMETIC_OPERATOR", p.peek())
		if aro {
			ado, adoerr := p.ParseAdditiveOperator()
			if adoerr != nil {
				return nil, roerr
			}
			se.Children = append(se.Children, ado)

			// expect <term>
			t, terr = p.ParseTerm()

			if terr != nil {
				return nil, terr
			}

			t.Children = append(t.Children, se)
			continue
		}

		break
	}

	return se, nil
}

// "<term>":                    {"<factor>", "(<multiplicative-operator> <factor>)*"},
func (p *Parser) ParseTerm() (*AbstractSyntaxTree, error) {
	// init + create main node
	t := NewNode("<term>")

	// expect <factor>
	f, ferr := p.ParseFactor()
	if ferr != nil {
		return nil, ferr
	}
	t.Children = append(t.Children, f)

	// accept (<multiplicative-operator> <factor>)*
	for {
		aro, _ := regexp.MatchString("^ARITHMETIC_OPERATOR", p.peek())
		if aro {
			ado, adoerr := p.ParseAdditiveOperator()
			if adoerr != nil {
				return nil, adoerr
			}
			t.Children = append(t.Children, ado)

			// expect <factor>
			f, ferr = p.ParseFactor()

			if ferr != nil {
				return nil, ferr
			}

			t.Children = append(t.Children, f)
			continue
		}

		break
	}

	return t, nil
}

// "<factor>":                  {"(IDENTIFIER | NUMBER | CHAR_LITERAL | STRING_LITERAL | ( LPARENTHESES(() <expression> RPARENTHESES()) ) | LOGICAL_OPERATOR(tidak))", "(<factor> | <function-declaration>)"},
func (p *Parser) ParseFactor() (*AbstractSyntaxTree, error) {
	// init + create main node
	f := NewNode("<factor>")

	// accept {"(IDENTIFIER | NUMBER | CHAR_LITERAL | STRING_LITERAL | ( LPARENTHESES(() <expression> RPARENTHESES()) ) | LOGICAL_OPERATOR(tidak))"
	for {
		i, _ := regexp.MatchString("^IDENTFIER", p.peek())
		n, _ := regexp.MatchString("^NUMBER", p.peek())
		cl, _ := regexp.MatchString("^CHAR_LITERAL", p.peek())
		sl, _ := regexp.MatchString("^STRING_LITERAL", p.peek())
		lp, _ := regexp.MatchString("^LPARENTHESES(()", p.peek())
		lo, _ := regexp.MatchString("^LOGICAL_OPERATOR", p.peek())

		if i {
			ierr := p.expect("IDENTFIER")
			if ierr != nil {
				return nil, ierr
			}
			f.Children = append(f.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		} else if n {
			nerr := p.expect("NUMBER")
			if nerr != nil {
				return nil, nerr
			}
			f.Children = append(f.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		} else if cl {
			clerr := p.expect("CHAR_LITERAL")
			if clerr != nil {
				return nil, clerr
			}
			f.Children = append(f.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		} else if sl {
			slerr := p.expect("STRING_LITERAL")
			if slerr != nil {
				return nil, slerr
			}
			f.Children = append(f.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		} else if lp {
			lperr := p.expect("LPARENTHESES(()")
			if lperr != nil {
				return nil, lperr
			}
			f.Children = append(f.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		} else if lo {
			loerr := p.expect("LOGICAL_OPERATOR")
			if loerr != nil {
				return nil, loerr
			}
			f.Children = append(f.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		}

		break
	}

	// accept (<factor> <function-declaration>)*
	for {
		i, _ := regexp.MatchString("^IDENTFIER", p.peek())
		n, _ := regexp.MatchString("^NUMBER", p.peek())
		cl, _ := regexp.MatchString("^CHAR_LITERAL", p.peek())
		sl, _ := regexp.MatchString("^STRING_LITERAL", p.peek())
		lp, _ := regexp.MatchString("^LPARENTHESES(()", p.peek())
		lo, _ := regexp.MatchString("^LOGICAL_OPERATOR", p.peek())

		if i || n || cl || sl || lp || lo {
			f, ferr := p.ParseFactor()
			if ferr != nil {
				return nil, ferr
			}
			f.Children = append(f.Children, f)
			continue
		} else if p.peek() == "KEYWORD(fungsi)" {
			fd, fderr := p.ParseFunctionDeclaration()
			if fderr != nil {
				return nil, fderr
			}
			f.Children = append(f.Children, fd)
			continue
		}
		break
	}

	return f, nil
}

// "<relational-operator>":     {"= | > | < | >= | <= | <>"},
func (p *Parser) ParseRelationalOperator() (*AbstractSyntaxTree, error) {
	// init + create main node
	ro := NewNode("<relational-operator>")

	for {
		robool, _ := regexp.MatchString("^RELATIONAL_OPERATOR", p.peek())

		if robool {
			roerr := p.expect("RELATIONAL_OPERATOR")
			if roerr != nil {
				return nil, roerr
			}
			ro.Children = append(ro.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		}
		break
	}

	return ro, nil
}

func (p *Parser) ParseAdditiveOperator() (*AbstractSyntaxTree, error) {
	// init + create main node
	ao := NewNode("<additive-operator>")

	for {
		aobool, _ := regexp.MatchString("^ARITHMETIC_OPERATOR", p.peek())

		if aobool {
			roerr := p.expect("ARITHMETIC_OPERATOR")
			if roerr != nil {
				return nil, roerr
			}
			ao.Children = append(ao.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		}
		break
	}

	return ao, nil
}

func (p *Parser) ParseMultiplicativeOperator() (*AbstractSyntaxTree, error) {
	// init + create main node
	mo := NewNode("<multiplicative-operator>")

	for {
		aobool, _ := regexp.MatchString("^ARITHMETIC_OPERATOR", p.peek())

		if aobool {
			roerr := p.expect("ARITHMETIC_OPERATOR")
			if roerr != nil {
				return nil, roerr
			}
			mo.Children = append(mo.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1]))
		}
		break
	}

	return mo, nil
}
