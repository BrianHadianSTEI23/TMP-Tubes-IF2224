package milestone3

import (
	"compiler/milestone2"
	"fmt"
	"strconv"
	"strings"
)

// SemanticAnalyzer performs semantic analysis on parse tree
// Builds symbol table and decorated AST simultaneously
type SemanticAnalyzer struct {
	SymTable      *SymbolTable
	CurrentOffset int
	Errors        []string
	Warnings      []string
}

// Create new semantic analyzer
func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{
		SymTable:      NewSymbolTable(),
		CurrentOffset: 5, // Stack frame header offset
		Errors:        make([]string, 0),
		Warnings:      make([]string, 0),
	}
}

// Analyze performs semantic analysis on parse tree
// Returns decorated AST and symbol table
func (sa *SemanticAnalyzer) Analyze(parseTree *milestone2.AbstractSyntaxTree) (DecoratedNode, error) {
	if parseTree == nil {
		return nil, fmt.Errorf("parse tree is nil")
	}

	// Visit program node - builds symbol table and decorated AST
	decoratedAST := sa.visitProgram(parseTree)

	// Return errors if any
	if len(sa.Errors) > 0 {
		return decoratedAST, fmt.Errorf("semantic analysis failed with %d error(s)", len(sa.Errors))
	}

	return decoratedAST, nil
}

// Get symbol table
func (sa *SemanticAnalyzer) GetSymbolTable() *SymbolTable {
	return sa.SymTable
}

// Get errors
func (sa *SemanticAnalyzer) GetErrors() []string {
	return sa.Errors
}

// Get warnings
func (sa *SemanticAnalyzer) GetWarnings() []string {
	return sa.Warnings
}

// ========== VISITOR FUNCTIONS ==========

// Visit <program> node
func (sa *SemanticAnalyzer) visitProgram(node *milestone2.AbstractSyntaxTree) *ProgramNode {
	if node.Value != "<program>" {
		sa.addError(fmt.Sprintf("expected <program> node, got %s", node.Value))
		return NewProgramNode("")
	}

	programName := ""
	var declarationsNode DecoratedNode
	var blockNode DecoratedNode

	// Extract program name from <program-header>
	if len(node.Children) > 0 {
		headerNode := node.Children[0]
		if headerNode.Value == "<program-header>" && len(headerNode.Children) > 1 {
			programNameNode := headerNode.Children[1]
			if strings.Contains(programNameNode.Value, "IDENTIFIER") {
				programName = extractValue(programNameNode.Value)

				// Add program to symbol table
				sa.SymTable.Enter(programName, ObjProgram, TypeNone, 0, 1, 0)
			}
		}
	}

	// Create program node
	programNode := NewProgramNode(programName)

	// Set symbol table info
	if idx, found := sa.SymTable.Lookup(programName); found {
		programNode.TabIndex = idx
		if entry, _ := sa.SymTable.GetEntry(idx); entry != nil {
			programNode.Level = entry.Lev
		}
	}

	// Visit <declaration-part> - builds symbol table + decorated AST
	if len(node.Children) > 1 {
		declarationsNode = sa.visitDeclarationPart(node.Children[1])
	}

	// Visit <compound-statement> - semantic validation + decorated AST
	if len(node.Children) > 2 && node.Children[2].Value == "<compound-statement>" {
		blockNode = sa.visitCompoundStatement(node.Children[2])
	}

	programNode.Declarations = declarationsNode
	programNode.Block = blockNode

	return programNode
}

// Visit <declaration-part> node
func (sa *SemanticAnalyzer) visitDeclarationPart(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	if node.Value != "<declaration-part>" {
		return nil
	}

	declarations := make([]DecoratedNode, 0)

	// Visit all declarations
	for _, child := range node.Children {
		var decl DecoratedNode
		switch child.Value {
		case "<const-declaration>":
			decl = sa.visitConstDeclaration(child)
		case "<var-declaration>":
			decl = sa.visitVarDeclaration(child)
		case "<type-declaration>":
			decl = sa.visitTypeDeclaration(child)
		case "<subprogram-declaration>":
			decl = sa.visitSubprogramDeclaration(child)
		}

		// Flatten DeclarationListNode to avoid nesting
		if decl != nil {
			if declList, ok := decl.(*DeclarationListNode); ok {
				declarations = append(declarations, declList.Declarations...)
			} else {
				declarations = append(declarations, decl)
			}
		}
	}

	if len(declarations) == 0 {
		return nil
	}

	if len(declarations) == 1 {
		return declarations[0]
	}

	return NewDeclarationListNode(declarations)
}

// Visit <var-declaration> node
// Builds symbol table AND creates VarDecl nodes
func (sa *SemanticAnalyzer) visitVarDeclaration(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	declarations := make([]DecoratedNode, 0)

	// Process variable declarations: <identifier-list> : <type> ;
	for i := 1; i < len(node.Children); i += 4 {
		if i+2 >= len(node.Children) {
			break
		}

		identifierListNode := node.Children[i]
		typeNode := node.Children[i+2]

		// Extract identifiers
		identifiers := sa.extractIdentifierList(identifierListNode)

		// Process type
		typ, ref := sa.processType(typeNode)
		typeSize := sa.SymTable.getTypeSize(typ, ref)

		// Enter each identifier to symbol table AND create decorated node
		for _, identifier := range identifiers {
			// Check for duplicate
			if sa.SymTable.IsDeclaredInCurrentScope(identifier) {
				sa.addError(fmt.Sprintf("Duplicate variable declaration: %s", identifier))
				continue
			}

			// Add to symbol table
			tabIndex := sa.SymTable.Enter(
				identifier,
				ObjVariable,
				typ,
				ref,
				1, // normal variable
				sa.CurrentOffset,
			)
			sa.SymTable.AddVariableSize(typeSize)
			sa.CurrentOffset += typeSize

			// Create decorated node
			varDecl := NewVarDeclNode(identifier, typ)
			varDecl.TabIndex = tabIndex
			varDecl.Type = typ
			varDecl.Ref = ref
			varDecl.Level = sa.SymTable.CurrentLevel
			varDecl.Address = sa.CurrentOffset - typeSize

			declarations = append(declarations, varDecl)
		}
	}

	if len(declarations) == 0 {
		return nil
	}

	if len(declarations) == 1 {
		return declarations[0]
	}

	return NewDeclarationListNode(declarations)
}

// Visit <const-declaration> node
func (sa *SemanticAnalyzer) visitConstDeclaration(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	declarations := make([]DecoratedNode, 0)

	// Process const definitions
	for i := 1; i < len(node.Children); i++ {
		child := node.Children[i]
		if child.Value == "<const-def>" && len(child.Children) >= 3 {
			identifierNode := child.Children[0]
			valueNode := child.Children[2]

			identifier := extractValue(identifierNode.Value)

			if sa.SymTable.IsDeclaredInCurrentScope(identifier) {
				sa.addError(fmt.Sprintf("Duplicate constant declaration: %s", identifier))
				continue
			}

			value, typ := sa.extractConstValue(valueNode.Value)

			// Add to symbol table
			tabIndex := sa.SymTable.Enter(identifier, ObjConstant, typ, -1, 1, value)

			// Create decorated node
			constDecl := NewConstDeclNode(identifier, value, typ)
			constDecl.TabIndex = tabIndex
			constDecl.Type = typ
			constDecl.Level = sa.SymTable.CurrentLevel

			declarations = append(declarations, constDecl)
		}
	}

	if len(declarations) == 0 {
		return nil
	}

	if len(declarations) == 1 {
		return declarations[0]
	}

	return NewDeclarationListNode(declarations)
}

// Visit <type-declaration> node
func (sa *SemanticAnalyzer) visitTypeDeclaration(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	// Type declarations don't appear in decorated AST typically
	// But we still process them for symbol table

	for i := 1; i < len(node.Children); i += 4 {
		if i+2 >= len(node.Children) {
			break
		}

		identifierNode := node.Children[i]
		typeNode := node.Children[i+2]

		identifier := extractValue(identifierNode.Value)

		if sa.SymTable.IsDeclaredInCurrentScope(identifier) {
			sa.addError(fmt.Sprintf("Duplicate type declaration: %s", identifier))
			continue
		}

		typ, ref := sa.processType(typeNode)
		sa.SymTable.Enter(identifier, ObjType, typ, ref, 1, 0)
	}

	return nil
}

// Visit <subprogram-declaration> node
func (sa *SemanticAnalyzer) visitSubprogramDeclaration(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	if len(node.Children) < 2 {
		return nil
	}

	keywordNode := node.Children[0]
	nameNode := node.Children[1]

	isFungsi := strings.Contains(keywordNode.Value, "fungsi")
	name := extractValue(nameNode.Value)

	if sa.SymTable.IsDeclaredInCurrentScope(name) {
		sa.addError(fmt.Sprintf("Duplicate subprogram declaration: %s", name))
		return nil
	}

	// Create new level and block
	blockIndex := sa.SymTable.enterLevelWithBlock()

	oldOffset := sa.CurrentOffset
	sa.CurrentOffset = 5

	// Process parameters
	parameters := make([]DecoratedNode, 0)
	lastParamIndex := 0

	for _, child := range node.Children {
		if child.Value == "<parameter-list>" {
			params := sa.extractParameters(child)
			for _, param := range params {
				size := sa.SymTable.getTypeSize(param.Type, param.Ref)
				lastParamIndex = sa.SymTable.Enter(param.Name, ObjVariable, param.Type, param.Ref, param.Nrm, sa.CurrentOffset)
				sa.SymTable.AddParameterSize(size)
				sa.CurrentOffset += size

				// Create parameter decorated node
				paramNode := NewVarDeclNode(param.Name, param.Type)
				paramNode.TabIndex = lastParamIndex
				paramNode.Type = param.Type
				paramNode.Level = sa.SymTable.CurrentLevel
				parameters = append(parameters, paramNode)
			}
			break
		}
	}

	sa.SymTable.UpdateBlockLastParam(blockIndex, lastParamIndex)

	// Determine return type
	returnType := TypeNone
	if isFungsi {
		for i := 0; i < len(node.Children)-1; i++ {
			if strings.Contains(node.Children[i].Value, "COLON") && i+1 < len(node.Children) {
				if node.Children[i+1].Value == "<type>" {
					returnType, _ = sa.processType(node.Children[i+1])
					break
				}
			}
		}
	}

	// Enter subprogram to parent scope
	objClass := ObjProcedure
	if isFungsi {
		objClass = ObjFunction
	}

	sa.SymTable.exitLevel()
	tabIndex := sa.SymTable.Enter(name, objClass, returnType, blockIndex, 1, 0)
	sa.SymTable.Display[sa.SymTable.CurrentLevel+1] = blockIndex
	sa.SymTable.enterLevel()

	// For functions, add implicit return variable
	if isFungsi {
		sa.SymTable.Enter(name, ObjVariable, returnType, -1, 1, 0)
	}

	// Process local declarations and body
	var localDecls DecoratedNode
	var body DecoratedNode

	for _, child := range node.Children {
		if child.Value == "<declaration-part>" {
			localDecls = sa.visitDeclarationPart(child)
		} else if child.Value == "<compound-statement>" {
			body = sa.visitCompoundStatement(child)
		}
	}

	sa.SymTable.exitLevel()
	sa.CurrentOffset = oldOffset

	// Create subprogram decorated node
	subprogNode := NewSubprogramDeclNode(name, parameters, returnType, body, isFungsi)
	subprogNode.TabIndex = tabIndex
	subprogNode.Type = returnType
	subprogNode.Level = sa.SymTable.CurrentLevel
	if localDecls != nil {
		// TODO: add local declarations to body
	}

	return subprogNode
}

// Visit <compound-statement> node
// Performs semantic validation while building decorated AST
func (sa *SemanticAnalyzer) visitCompoundStatement(node *milestone2.AbstractSyntaxTree) *BlockNode {
	statements := make([]DecoratedNode, 0)

	// Find <statement-list>
	for _, child := range node.Children {
		if child.Value == "<statement-list>" {
			statements = sa.visitStatementList(child)
			break
		}
	}

	blockNode := NewBlockNode(statements)
	blockNode.BlockIndex = sa.SymTable.CurrentBlock
	blockNode.Level = sa.SymTable.CurrentLevel

	return blockNode
}

// Visit <statement-list> node
func (sa *SemanticAnalyzer) visitStatementList(node *milestone2.AbstractSyntaxTree) []DecoratedNode {
	statements := make([]DecoratedNode, 0)

	for _, child := range node.Children {
		if stmt := sa.visitStatement(child); stmt != nil {
			statements = append(statements, stmt)
		}
	}

	return statements
}

// Visit individual statement
func (sa *SemanticAnalyzer) visitStatement(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	switch node.Value {
	case "<assignment-statement>":
		return sa.visitAssignmentStatement(node)
	case "<procedure-call>":
		return sa.visitProcedureCall(node)
	case "<compound-statement>":
		return sa.visitCompoundStatement(node)
	case "<if-statement>":
		return sa.visitIfStatement(node)
	case "<while-statement>":
		return sa.visitWhileStatement(node)
	case "<for-statement>":
		return sa.visitForStatement(node)
	default:
		// Skip non-statement nodes like SEMICOLON
		return nil
	}
}

// Visit <assignment-statement> node
// Semantic rule: assignment_statement.node = new AssignNode(new VarNode(ID.lexeme), expr.node)
func (sa *SemanticAnalyzer) visitAssignmentStatement(node *milestone2.AbstractSyntaxTree) *AssignNode {
	var targetName string
	var valueNode DecoratedNode

	// Extract target variable and expression
	for _, child := range node.Children {
		if child.Value == "<variable>" {
			// Extract identifier from variable node
			for _, grandchild := range child.Children {
				if strings.Contains(grandchild.Value, "IDENTIFIER") {
					targetName = extractValue(grandchild.Value)
					break
				}
			}
		} else if child.Value == "<expression>" {
			valueNode = sa.visitExpression(child)
		}
	}

	// Look up target variable in symbol table
	var targetNode *VarNode
	if targetName != "" {
		tabIndex, found := sa.SymTable.Lookup(targetName)
		if !found {
			sa.addError(fmt.Sprintf("Undefined variable '%s'", targetName))
			targetNode = NewVarNode(targetName)
		} else {
			entry, _ := sa.SymTable.GetEntry(tabIndex)
			if entry != nil {
				if entry.Obj != ObjVariable {
					sa.addError(fmt.Sprintf("'%s' is not a variable", targetName))
				}

				targetNode = NewVarNode(targetName)
				targetNode.TabIndex = tabIndex
				targetNode.Type = entry.Type
				targetNode.Ref = entry.Ref
				targetNode.Level = entry.Lev
				targetNode.Address = entry.Adr
				targetNode.IsLValue = true

				// Type checking
				if valueNode != nil {
					targetType := entry.Type
					valueType := sa.getNodeType(valueNode)

					if !sa.typesCompatible(targetType, valueType) {
						sa.addError(fmt.Sprintf("Type mismatch in assignment: cannot assign %s to %s", valueType, targetType))
					}
				}
			}
		}
	}

	if targetNode == nil {
		targetNode = NewVarNode("error")
	}

	if valueNode == nil {
		valueNode = NewNumberNode(0)
	}

	return NewAssignNode(targetNode, valueNode)
}

// Visit <expression> node
// Semantic rules:
// - <expression> → <simple-expr> : expression.node = simple_expr.node
// - <expression> → <simple-expr> relop <simple-expr> : expression.node = new BinOpNode(relop, left, right)
func (sa *SemanticAnalyzer) visitExpression(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	if len(node.Children) == 1 {
		// expression → simple-expression
		return sa.visitSimpleExpression(node.Children[0])
	}

	if len(node.Children) == 3 {
		// expression → simple-expression relop simple-expression
		left := sa.visitSimpleExpression(node.Children[0])
		operator := extractValue(node.Children[1].Value)
		right := sa.visitSimpleExpression(node.Children[2])

		// Type checking
		leftType := sa.getNodeType(left)
		rightType := sa.getNodeType(right)

		if !sa.typesCompatible(leftType, rightType) {
			sa.addError(fmt.Sprintf("Type mismatch in relational operation: %s and %s", leftType, rightType))
		}

		binOp := NewBinOpNode(operator, left, right)
		binOp.Type = TypeBoolean // Relational operators return boolean
		return binOp
	}

	return NewNumberNode(0)
}

// Visit <simple-expression> node
// Semantic rule: <simple-expr> → <term> (addop <term>)*
func (sa *SemanticAnalyzer) visitSimpleExpression(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	if len(node.Children) == 1 {
		// simple-expression → term
		return sa.visitTerm(node.Children[0])
	}

	if len(node.Children) >= 3 {
		// simple-expression → term addop term
		left := sa.visitTerm(node.Children[0])
		operator := extractValue(node.Children[1].Value)
		right := sa.visitTerm(node.Children[2])

		// Type checking
		leftType := sa.getNodeType(left)
		rightType := sa.getNodeType(right)

		if !sa.isNumericType(leftType) || !sa.isNumericType(rightType) {
			sa.addError(fmt.Sprintf("Arithmetic operator requires numeric operands"))
		}

		binOp := NewBinOpNode(operator, left, right)
		binOp.Type = TypeInteger
		return binOp
	}

	return NewNumberNode(0)
}

// Visit <term> node
// Semantic rule: <term> → <factor> (mulop <factor>)*
func (sa *SemanticAnalyzer) visitTerm(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	if len(node.Children) == 0 {
		return NewNumberNode(0)
	}

	// Start with first factor
	result := sa.visitFactor(node.Children[0])

	// Handle chained operators: factor (*|/|div|mod|dan) factor (*|/|div|mod|dan) factor ...
	for i := 1; i < len(node.Children)-1; i += 2 {
		if i+1 >= len(node.Children) {
			break
		}

		operatorNode := node.Children[i]
		operator := extractValue(operatorNode.Value)

		if node.Children[i+1].Value != "<factor>" {
			break
		}

		right := sa.visitFactor(node.Children[i+1])

		// Type checking
		leftType := sa.getNodeType(result)
		rightType := sa.getNodeType(right)

		if operator == "dan" || operator == "and" {
			// Logical AND - expects boolean operands
			if leftType != TypeBoolean || rightType != TypeBoolean {
				sa.addError("Logical AND operator requires boolean operands")
			}
			binOp := NewBinOpNode(operator, result, right)
			binOp.Type = TypeBoolean
			result = binOp
		} else {
			// Arithmetic operators
			if !sa.isNumericType(leftType) || !sa.isNumericType(rightType) {
				sa.addError("Multiplicative operator requires numeric operands")
			}
			binOp := NewBinOpNode(operator, result, right)
			binOp.Type = TypeInteger
			result = binOp
		}
	}

	return result
}

// Visit <factor> node
// Semantic rules:
// - <factor> → NUMBER : factor.node = new NumberNode(NUMBER.value)
// - <factor> → ID : factor.node = new VarNode(ID.lexeme)
// - <factor> → STRING_LITERAL : factor.node = new StringNode(STRING_LITERAL.value)
// - <factor> → true/false : factor.node = new BooleanNode(value)
// - <factor> → ( <expression> ) : factor.node = expression.node
// - <factor> → tidak <factor> : factor.node = new UnaryOpNode("tidak", factor.node)
// - <factor> → <function-call> : factor.node = function_call.node
// - <factor> → CHAR_LITERAL : factor.node = new CharNode(value)
func (sa *SemanticAnalyzer) visitFactor(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	for _, child := range node.Children {
		if strings.Contains(child.Value, "NUMBER") {
			// factor → NUMBER
			valueStr := extractValue(child.Value)
			value, _ := strconv.Atoi(valueStr)
			numberNode := NewNumberNode(value)
			numberNode.Type = TypeInteger
			return numberNode
		} else if strings.Contains(child.Value, "STRING_LITERAL") {
			// factor → STRING_LITERAL
			valueStr := extractValue(child.Value)
			// Remove quotes
			if len(valueStr) >= 2 && valueStr[0] == '\'' {
				valueStr = valueStr[1 : len(valueStr)-1]
			}
			stringNode := NewStringNode(valueStr)
			stringNode.Type = TypeChar
			return stringNode
		} else if strings.Contains(child.Value, "CHAR_LITERAL") {
			// factor → CHAR_LITERAL
			valueStr := extractValue(child.Value)
			if len(valueStr) >= 2 && valueStr[0] == '\'' {
				valueStr = valueStr[1 : len(valueStr)-1]
			}
			charVal := rune(0)
			if len(valueStr) > 0 {
				charVal = rune(valueStr[0])
			}
			charNode := &CharNode{
				BaseDecoratedNode: BaseDecoratedNode{Type: TypeChar},
				Value:             charVal,
			}
			return charNode
		} else if strings.Contains(child.Value, "true") || strings.Contains(child.Value, "false") {
			// factor → true | false
			valueStr := extractValue(child.Value)
			boolVal := (valueStr == "true")
			boolNode := NewBooleanNode(boolVal)
			boolNode.Type = TypeBoolean
			return boolNode
		} else if strings.Contains(child.Value, "tidak") || strings.Contains(child.Value, "not") {
			// factor → tidak <factor> (NOT operator)
			// Find the factor child
			for _, subChild := range node.Children {
				if subChild.Value == "<factor>" {
					operand := sa.visitFactor(subChild)
					operandType := sa.getNodeType(operand)
					if operandType != TypeBoolean {
						sa.addError("NOT operator requires boolean operand")
					}
					unaryOp := &UnaryOpNode{
						BaseDecoratedNode: BaseDecoratedNode{Type: TypeBoolean},
						Operator:          "tidak",
						Operand:           operand,
					}
					return unaryOp
				}
			}
		} else if child.Value == "<function-call>" {
			// factor → <function-call>
			return sa.visitFunctionCall(child)
		} else if child.Value == "<variable>" {
			// factor → variable (which contains ID)
			return sa.visitVariable(child)
		} else if strings.Contains(child.Value, "IDENTIFIER") {
			// factor → ID (direct identifier)
			return sa.visitIdentifier(child)
		} else if child.Value == "<expression>" {
			// factor → ( expression )
			return sa.visitExpression(child)
		} else {
			// Recursively check children
			if result := sa.visitFactor(child); result != nil {
				return result
			}
		}
	}

	return NewNumberNode(0)
}

// Visit <variable> node
// Handles:
// - <variable> → ID : simple variable reference
// - <variable> → ID [ expression ] : array indexing
// - <variable> → ID . ID : record field access
func (sa *SemanticAnalyzer) visitVariable(node *milestone2.AbstractSyntaxTree) DecoratedNode {
	var varNode *VarNode
	var indexExpr DecoratedNode
	var fieldName string

	// Extract identifier, optional index, and optional field name
	for _, child := range node.Children {
		if strings.Contains(child.Value, "IDENTIFIER") {
			if varNode == nil {
				// First identifier is the variable
				varNode = sa.visitIdentifier(child)
			} else {
				// Second identifier is field name (record.field)
				fieldName = extractValue(child.Value)
			}
		} else if child.Value == "<expression>" {
			// Array indexing: ID [ expression ]
			indexExpr = sa.visitExpression(child)
		}
	}

	if varNode == nil {
		return NewVarNode("unknown")
	}

	// Handle record field access first (ID.field)
	if fieldName != "" {
		// Verify variable is a record
		if varNode.Type != TypeRecord {
			sa.addError(fmt.Sprintf("'%s' is not a record", varNode.Name))
			return varNode
		}

		// Look up field in record's BTAB
		if varNode.Ref >= 0 && varNode.Ref < len(sa.SymTable.Btab) {
			btabEntry := sa.SymTable.Btab[varNode.Ref]

			// Search for field in record's symbol table entries
			fieldTabIndex := btabEntry.Last
			found := false

			for fieldTabIndex >= 0 && fieldTabIndex < len(sa.SymTable.Tab) {
				fieldEntry := sa.SymTable.Tab[fieldTabIndex]

				if fieldEntry.Identifier == fieldName {
					// Found the field - create new VarNode with field info
					fieldNode := NewVarNode(varNode.Name + "." + fieldName)
					fieldNode.TabIndex = fieldTabIndex
					fieldNode.Type = fieldEntry.Type
					fieldNode.Ref = fieldEntry.Ref
					fieldNode.Level = fieldEntry.Lev
					fieldNode.Address = varNode.Address + fieldEntry.Adr // Base + field offset
					fieldNode.IsLValue = true
					varNode = fieldNode
					found = true
					break
				}

				// Follow linked list
				fieldTabIndex = fieldEntry.Link
				if fieldTabIndex == -1 {
					break
				}
			}

			if !found {
				sa.addError(fmt.Sprintf("Record '%s' has no field '%s'", varNode.Name, fieldName))
			}
		}
	}

	// Handle array indexing (can be combined with record: record.field[index])
	if indexExpr != nil {
		varNode.IsIndexed = true
		varNode.Index = indexExpr

		// Type checking: verify variable is an array
		if varNode.Type != TypeArray {
			sa.addError(fmt.Sprintf("'%s' is not an array", varNode.Name))
		} else {
			// Get element type from ATAB using Ref
			if varNode.Ref >= 0 && varNode.Ref < len(sa.SymTable.Atab) {
				atabEntry := sa.SymTable.Atab[varNode.Ref]
				// Resolve to element type
				varNode.Type = TypeKind(atabEntry.Etyp)
				varNode.Ref = atabEntry.Eref

				// Type check index expression
				indexType := sa.getNodeType(indexExpr)
				if indexType != TypeInteger {
					sa.addError("Array index must be integer type")
				}
			}
		}
	}

	return varNode
}

// Visit IDENTIFIER node
// Semantic rule: Creates VarNode with symbol table lookup
func (sa *SemanticAnalyzer) visitIdentifier(node *milestone2.AbstractSyntaxTree) *VarNode {
	varName := extractValue(node.Value)
	varNode := NewVarNode(varName)

	// Look up in symbol table
	tabIndex, found := sa.SymTable.Lookup(varName)
	if !found {
		sa.addError(fmt.Sprintf("Undefined identifier '%s'", varName))
		return varNode
	}

	entry, _ := sa.SymTable.GetEntry(tabIndex)
	if entry != nil {
		varNode.TabIndex = tabIndex
		varNode.Type = entry.Type
		varNode.Ref = entry.Ref
		varNode.Level = entry.Lev
		varNode.Address = entry.Adr
		varNode.IsLValue = (entry.Obj == ObjVariable)
	}

	return varNode
}

// Visit <procedure-call> node
// Semantic rule: procedure_call.node = new ProcCallNode(ID.lexeme, params.nodes)
func (sa *SemanticAnalyzer) visitProcedureCall(node *milestone2.AbstractSyntaxTree) *ProcCallNode {
	var procName string
	arguments := make([]DecoratedNode, 0)

	// Extract procedure name and arguments
	for _, child := range node.Children {
		if strings.Contains(child.Value, "IDENTIFIER") || strings.Contains(child.Value, "writeln") || strings.Contains(child.Value, "write") {
			if strings.Contains(child.Value, "IDENTIFIER") {
				procName = extractValue(child.Value)
			} else {
				procName = extractValue(child.Value)
			}
		} else if child.Value == "<parameter-list>" {
			// Extract arguments from parameter list
			for _, paramChild := range child.Children {
				if paramChild.Value == "<expression>" {
					if arg := sa.visitExpression(paramChild); arg != nil {
						arguments = append(arguments, arg)
					}
				}
			}
		}
	}

	procCall := NewProcCallNode(procName, arguments)

	// Check if built-in
	builtIns := map[string]bool{
		"write": true, "writeln": true, "read": true, "readln": true,
	}

	if builtIns[procName] {
		procCall.IsBuiltIn = true
		procCall.TabIndex = -1
	} else {
		// Look up in symbol table
		if tabIndex, found := sa.SymTable.Lookup(procName); found {
			entry, _ := sa.SymTable.GetEntry(tabIndex)
			if entry != nil {
				if entry.Obj != ObjProcedure && entry.Obj != ObjFunction {
					sa.addError(fmt.Sprintf("'%s' is not a procedure", procName))
				}
				procCall.TabIndex = tabIndex
				procCall.Type = entry.Type
			}
		} else {
			sa.addError(fmt.Sprintf("Undefined procedure '%s'", procName))
		}
	}

	return procCall
}

// Visit <function-call> node (for function calls in expressions)
// Semantic rule: function_call.node = new ProcCallNode(ID.lexeme, params.nodes)
func (sa *SemanticAnalyzer) visitFunctionCall(node *milestone2.AbstractSyntaxTree) *ProcCallNode {
	var funcName string
	arguments := make([]DecoratedNode, 0)

	// Extract function name and arguments
	for _, child := range node.Children {
		if strings.Contains(child.Value, "IDENTIFIER") {
			funcName = extractValue(child.Value)
		} else if child.Value == "<expr-list>" {
			// Process expression list
			for _, exprChild := range child.Children {
				if exprChild.Value == "<expression>" {
					arg := sa.visitExpression(exprChild)
					arguments = append(arguments, arg)
				}
			}
		}
	}

	funcCall := NewProcCallNode(funcName, arguments)

	// Look up in symbol table
	tabIndex, found := sa.SymTable.Lookup(funcName)
	if !found {
		sa.addError(fmt.Sprintf("Undefined function '%s'", funcName))
	} else {
		entry, _ := sa.SymTable.GetEntry(tabIndex)
		if entry != nil {
			if entry.Obj != ObjFunction {
				sa.addError(fmt.Sprintf("'%s' is not a function", funcName))
			}
			funcCall.TabIndex = tabIndex
			funcCall.Type = entry.Type // Function return type
		}
	}

	return funcCall
}

// Visit <if-statement> node
func (sa *SemanticAnalyzer) visitIfStatement(node *milestone2.AbstractSyntaxTree) *IfNode {
	ifNode := &IfNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
	}

	// Extract condition and statements
	for _, child := range node.Children {
		if child.Value == "<expression>" && ifNode.Condition == nil {
			ifNode.Condition = sa.visitExpression(child)

			// Type check: condition must be boolean
			condType := sa.getNodeType(ifNode.Condition)
			if condType != TypeBoolean {
				sa.addError("If condition must be boolean type")
			}
		} else if child.Value == "<statement>" || child.Value == "<compound-statement>" {
			if ifNode.ThenStmt == nil {
				ifNode.ThenStmt = sa.visitStatement(child)
			} else if ifNode.ElseStmt == nil {
				ifNode.ElseStmt = sa.visitStatement(child)
			}
		}
	}

	return ifNode
}

// Visit <while-statement> node
func (sa *SemanticAnalyzer) visitWhileStatement(node *milestone2.AbstractSyntaxTree) *WhileNode {
	whileNode := &WhileNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
	}

	// Extract condition and body
	for _, child := range node.Children {
		if child.Value == "<expression>" {
			whileNode.Condition = sa.visitExpression(child)

			// Type check: condition must be boolean
			condType := sa.getNodeType(whileNode.Condition)
			if condType != TypeBoolean {
				sa.addError("While condition must be boolean type")
			}
		} else if child.Value == "<statement>" || child.Value == "<compound-statement>" {
			whileNode.Body = sa.visitStatement(child)
		}
	}

	return whileNode
}

// Visit <for-statement> node
func (sa *SemanticAnalyzer) visitForStatement(node *milestone2.AbstractSyntaxTree) *ForNode {
	forNode := &ForNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
	}

	// Extract for loop components
	// TODO: Implement for statement parsing

	return forNode
}

// ========== HELPER FUNCTIONS ==========

// Process type node
func (sa *SemanticAnalyzer) processType(node *milestone2.AbstractSyntaxTree) (TypeKind, int) {
	if node.Value == "<array-type>" {
		return sa.processArrayType(node)
	}

	if node.Value == "<record-type>" {
		return sa.processRecordType(node)
	}

	if node.Value != "<type>" {
		return TypeNone, -1
	}

	if len(node.Children) == 0 {
		return TypeNone, -1
	}

	child := node.Children[0]

	if strings.Contains(child.Value, "KEYWORD") {
		typeStr := extractValue(child.Value)
		switch typeStr {
		case "integer":
			return TypeInteger, -1
		case "boolean":
			return TypeBoolean, -1
		case "char":
			return TypeChar, -1
		}
	}

	if child.Value == "<array-type>" {
		return sa.processArrayType(child)
	}

	if child.Value == "<record-type>" {
		return sa.processRecordType(child)
	}

	if strings.Contains(child.Value, "IDENTIFIER") {
		typeName := extractValue(child.Value)
		idx, found := sa.SymTable.Lookup(typeName)
		if found && sa.SymTable.Tab[idx].Obj == ObjType {
			return sa.SymTable.Tab[idx].Type, sa.SymTable.Tab[idx].Ref
		}
		sa.addError(fmt.Sprintf("Undefined type '%s'", typeName))
	}

	return TypeNone, -1
}

// Process array type
func (sa *SemanticAnalyzer) processArrayType(node *milestone2.AbstractSyntaxTree) (TypeKind, int) {
	low, high := 0, 0
	var elementTypeNode *milestone2.AbstractSyntaxTree

	for i, child := range node.Children {
		if strings.Contains(child.Value, "NUMBER") {
			if low == 0 {
				low, _ = strconv.Atoi(extractValue(child.Value))
			} else {
				high, _ = strconv.Atoi(extractValue(child.Value))
			}
		} else if child.Value == "<type>" {
			elementTypeNode = child
		} else if i == len(node.Children)-1 {
			elementTypeNode = child
		}
	}

	if elementTypeNode == nil {
		return TypeArray, -1
	}

	elemType, elemRef := sa.processType(elementTypeNode)
	elemSize := sa.SymTable.getTypeSize(elemType, elemRef)

	atabIndex := sa.SymTable.EnterArray(0, int(elemType), elemRef, low, high, elemSize)

	return TypeArray, atabIndex
}

// Process record type
// Creates BTAB entry for record and processes field declarations
func (sa *SemanticAnalyzer) processRecordType(node *milestone2.AbstractSyntaxTree) (TypeKind, int) {
	// Create new block for record type
	oldBlock := sa.SymTable.CurrentBlock
	oldOffset := sa.CurrentOffset
	sa.CurrentOffset = 0 // Fields start at offset 0

	blockIndex := sa.SymTable.enterBlock()

	// Process field-list
	for _, child := range node.Children {
		if child.Value == "<field-list>" {
			sa.processFieldList(child, blockIndex)
			break
		}
	}

	// Update block's vsze with total size of all fields
	sa.SymTable.Btab[blockIndex].Vsze = sa.CurrentOffset
	sa.SymTable.Btab[blockIndex].Lpar = 0 // Records have no parameters

	// Restore state
	sa.SymTable.CurrentBlock = oldBlock
	sa.CurrentOffset = oldOffset

	return TypeRecord, blockIndex
}

// Process field list for record type
func (sa *SemanticAnalyzer) processFieldList(node *milestone2.AbstractSyntaxTree, blockIndex int) {
	if node.Value != "<field-list>" {
		return
	}

	// Parse pattern: <identifier-list> : <type> (; <identifier-list> : <type>)*
	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]

		if child.Value == "<identifier-list>" && i+2 < len(node.Children) {
			// Extract field names
			fieldNames := sa.extractIdentifierList(child)

			// Get type (skip colon at i+1, type at i+2)
			typeNode := node.Children[i+2]
			fieldType, fieldRef := sa.processType(typeNode)
			fieldSize := sa.SymTable.getTypeSize(fieldType, fieldRef)

			// Enter each field into symbol table
			for _, fieldName := range fieldNames {
				// Enter field with ObjField class
				tabIndex := sa.SymTable.Enter(
					fieldName,
					ObjField,
					fieldType,
					fieldRef,
					1, // nrm = 1 for fields (normal)
					sa.CurrentOffset,
				)

				// Update last pointer in BTAB
				if blockIndex >= 0 && blockIndex < len(sa.SymTable.Btab) {
					if sa.SymTable.Btab[blockIndex].Last == -1 {
						sa.SymTable.Btab[blockIndex].Last = tabIndex
					}
				}

				sa.CurrentOffset += fieldSize
			}

			// Skip to next field group (skip colon and type)
			i += 2
		}
	}
}

// Extract identifier list
func (sa *SemanticAnalyzer) extractIdentifierList(node *milestone2.AbstractSyntaxTree) []string {
	identifiers := make([]string, 0)

	for _, child := range node.Children {
		if strings.Contains(child.Value, "IDENTIFIER") {
			identifiers = append(identifiers, extractValue(child.Value))
		}
	}

	return identifiers
}

// Parameter struct
type Parameter struct {
	Name string
	Type TypeKind
	Ref  int
	Nrm  int
}

// Extract parameters
// Parses <parameter-list> and extracts parameter info including nrm field
// nrm = 1 for normal (value) parameters
// nrm = 0 for var (reference) parameters
func (sa *SemanticAnalyzer) extractParameters(node *milestone2.AbstractSyntaxTree) []Parameter {
	params := make([]Parameter, 0)

	if node.Value != "<parameter-list>" {
		return params
	}

	// Parse pattern: <identifier-list> : <type> (; <identifier-list> : <type>)*
	// Note: Current parser doesn't distinguish var parameters
	// If parser is extended to support "variabel" keyword before identifier-list,
	// set nrm=0 for those parameters

	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]

		// Check if this is a var parameter (variabel keyword)
		isVarParam := false
		if strings.Contains(child.Value, "variabel") {
			isVarParam = true
			i++ // Skip to identifier-list
			if i >= len(node.Children) {
				break
			}
			child = node.Children[i]
		}

		if child.Value == "<identifier-list>" && i+2 < len(node.Children) {
			// Extract parameter names
			paramNames := sa.extractIdentifierList(child)

			// Get type (skip colon at i+1, type at i+2)
			typeNode := node.Children[i+2]
			paramType, paramRef := sa.processType(typeNode)

			// Determine nrm value
			nrm := 1 // Default: normal (value) parameter
			if isVarParam {
				nrm = 0 // Var (reference) parameter
			}

			// Create parameter entries
			for _, paramName := range paramNames {
				params = append(params, Parameter{
					Name: paramName,
					Type: paramType,
					Ref:  paramRef,
					Nrm:  nrm,
				})
			}

			// Skip to next parameter group (skip colon and type)
			i += 2
		}
	}

	return params
}

// Extract constant value
func (sa *SemanticAnalyzer) extractConstValue(tokenValue string) (int, TypeKind) {
	if strings.Contains(tokenValue, "NUMBER") {
		valueStr := extractValue(tokenValue)
		value, _ := strconv.Atoi(valueStr)
		return value, TypeInteger
	} else if strings.Contains(tokenValue, "true") || strings.Contains(tokenValue, "false") {
		if strings.Contains(tokenValue, "true") {
			return 1, TypeBoolean
		}
		return 0, TypeBoolean
	} else if strings.Contains(tokenValue, "STRING_LITERAL") {
		return 0, TypeChar
	}
	return 0, TypeNone
}

// Get node type
func (sa *SemanticAnalyzer) getNodeType(node DecoratedNode) TypeKind {
	switch n := node.(type) {
	case *VarNode:
		return n.Type
	case *NumberNode:
		return TypeInteger
	case *StringNode:
		return TypeChar
	case *BooleanNode:
		return TypeBoolean
	case *BinOpNode:
		return n.Type
	case *UnaryOpNode:
		return n.Type
	default:
		return TypeNone
	}
}

// Check type compatibility
func (sa *SemanticAnalyzer) typesCompatible(type1, type2 TypeKind) bool {
	if type1 == type2 {
		return true
	}
	// Integer and char are compatible
	if (type1 == TypeInteger && type2 == TypeChar) ||
		(type1 == TypeChar && type2 == TypeInteger) {
		return true
	}
	return false
}

// Check if type is numeric
func (sa *SemanticAnalyzer) isNumericType(typ TypeKind) bool {
	return typ == TypeInteger || typ == TypeChar
}

// Add error
func (sa *SemanticAnalyzer) addError(message string) {
	sa.Errors = append(sa.Errors, message)
}

// Add warning
func (sa *SemanticAnalyzer) addWarning(message string) {
	sa.Warnings = append(sa.Warnings, message)
}

// Extract value from token
func extractValue(tokenValue string) string {
	if start := strings.Index(tokenValue, "("); start >= 0 {
		if end := strings.LastIndex(tokenValue, ")"); end >= 0 {
			return tokenValue[start+1 : end]
		}
	}
	return tokenValue
}
