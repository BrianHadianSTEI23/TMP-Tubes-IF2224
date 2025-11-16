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
	m, _ := regexp.MatchString(`^`+tok, tok)
	if m {
		p.next()
		return true
	}
	return false
}
func (p *Parser) expect(tok string) error {
	m, _ := regexp.MatchString(`^`+tok, tok)
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
	var acceptValue bool

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
			acceptValue = p.accept("NUMBER")
			if nerr != nil {
				return nil, nerr
			}
			td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable
		} else if c {
			acceptValue = p.accept("CHAR_LITERAL")
			if cerr != nil {
				return nil, cerr
			}
			td.Children = append(td.Children, NewLeaf(p.Tokens[p.Pos-1], p.Tokens[p.Pos-1])) // this is minus one because p.expect has incremented the Pos variable
		} else if s {
			acceptValue = p.accept("STRING_LITERAL")
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
	var err error

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

// "<function-declaration>":    {"KEYWORD(function)", "IDENTIFIER", "(formal-parameter-list)*", "SEMICOLON(;)"},
func (p *Parser) ParseFunctionDeclaration() (*AbstractSyntaxTree, error) {
	// init + create main node
	fd := NewNode("<function-declaration>")
	var err error

	// expect KEYWORD(function)
	err = p.expect("KEYWORD(function)")
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
func (p *Parser) ParseStatement() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseStatementList() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseAssignmentStatement() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseIfStatement() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseWhileStatement() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseForStatement() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseParameterList() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseExpression() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseSimpleExpression() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseTerm() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseFactor() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseRelationalOperator() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseAdditiveOperator() (*AbstractSyntaxTree, error) {
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

func (p *Parser) ParseMultiplicativeOperator() (*AbstractSyntaxTree, error) {
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
