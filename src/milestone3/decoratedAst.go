package milestone3

import "fmt"

// DecoratedNode interface - base interface for all decorated AST nodes
type DecoratedNode interface {
	GetTabIndex() int
	GetType() TypeKind
	GetLevel() int
	Accept(visitor DecoratedNodeVisitor)
}

// DecoratedNodeVisitor interface for visitor pattern
type DecoratedNodeVisitor interface {
	VisitProgram(*ProgramNode)
	VisitVarDecl(*VarDeclNode)
	VisitConstDecl(*ConstDeclNode)
	VisitAssign(*AssignNode)
	VisitBinOp(*BinOpNode)
	VisitUnaryOp(*UnaryOpNode)
	VisitVar(*VarNode)
	VisitNumber(*NumberNode)
	VisitString(*StringNode)
	VisitBoolean(*BooleanNode)
	VisitBlock(*BlockNode)
	VisitProcCall(*ProcCallNode)
	VisitIf(*IfNode)
	VisitWhile(*WhileNode)
	VisitFor(*ForNode)
}

// BaseDecoratedNode - common fields for all decorated nodes
type BaseDecoratedNode struct {
	TabIndex int      // Index in symbol table
	Type     TypeKind // Type of the node
	Ref      int      // Reference to ATAB or BTAB (-1 if none)
	Level    int      // Lexical level
	Address  int      // Memory address/offset
	Errors   []string // Semantic errors for this node
	Warnings []string // Semantic warnings for this node
}

func (n *BaseDecoratedNode) GetTabIndex() int {
	return n.TabIndex
}

func (n *BaseDecoratedNode) GetType() TypeKind {
	return n.Type
}

func (n *BaseDecoratedNode) GetLevel() int {
	return n.Level
}

// ProgramNode - represents the entire program
type ProgramNode struct {
	BaseDecoratedNode
	Name         string
	Declarations DecoratedNode
	Block        DecoratedNode
}

func NewProgramNode(name string) *ProgramNode {
	return &ProgramNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Name: name,
	}
}

func (n *ProgramNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitProgram(n)
}

// DeclarationListNode - list of declarations
type DeclarationListNode struct {
	BaseDecoratedNode
	Declarations []DecoratedNode
}

func NewDeclarationListNode(declarations []DecoratedNode) *DeclarationListNode {
	return &DeclarationListNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Declarations: declarations,
	}
}

func (n *DeclarationListNode) Accept(visitor DecoratedNodeVisitor) {
	// Visit all declarations
	for _, decl := range n.Declarations {
		if decl != nil {
			decl.Accept(visitor)
		}
	}
}

// VarDeclNode - variable declaration
type VarDeclNode struct {
	BaseDecoratedNode
	Name string
}

func NewVarDeclNode(name string, typ TypeKind) *VarDeclNode {
	return &VarDeclNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     typ,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Name: name,
	}
}

func (n *VarDeclNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitVarDecl(n)
}

// ConstDeclNode - constant declaration
type ConstDeclNode struct {
	BaseDecoratedNode
	Name  string
	Value interface{}
}

func NewConstDeclNode(name string, value interface{}, typ TypeKind) *ConstDeclNode {
	return &ConstDeclNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     typ,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Name:  name,
		Value: value,
	}
}

func (n *ConstDeclNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitConstDecl(n)
}

// TypeDeclNode - type declaration
type TypeDeclNode struct {
	BaseDecoratedNode
	Name string
}

// SubprogramDeclNode - procedure or function declaration
type SubprogramDeclNode struct {
	BaseDecoratedNode
	Name       string
	Parameters []DecoratedNode
	ReturnType TypeKind
	Body       DecoratedNode
	IsFunction bool
}

func NewSubprogramDeclNode(name string, parameters []DecoratedNode, returnType TypeKind, body DecoratedNode, isFunction bool) *SubprogramDeclNode {
	return &SubprogramDeclNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     returnType,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Name:       name,
		Parameters: parameters,
		ReturnType: returnType,
		Body:       body,
		IsFunction: isFunction,
	}
}

func (n *SubprogramDeclNode) Accept(visitor DecoratedNodeVisitor) {
	// TODO: Implement visitor
}

// BlockNode - compound statement (block of statements)
type BlockNode struct {
	BaseDecoratedNode
	Statements []DecoratedNode
	BlockIndex int // Index in BTAB
}

func NewBlockNode(statements []DecoratedNode) *BlockNode {
	return &BlockNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Statements: statements,
		BlockIndex: -1,
	}
}

func (n *BlockNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitBlock(n)
}

// AssignNode - assignment statement
type AssignNode struct {
	BaseDecoratedNode
	Target DecoratedNode
	Value  DecoratedNode
}

func NewAssignNode(target, value DecoratedNode) *AssignNode {
	return &AssignNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Target: target,
		Value:  value,
	}
}

func (n *AssignNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitAssign(n)
}

// BinOpNode - binary operation
type BinOpNode struct {
	BaseDecoratedNode
	Operator string
	Left     DecoratedNode
	Right    DecoratedNode
}

func NewBinOpNode(operator string, left, right DecoratedNode) *BinOpNode {
	return &BinOpNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Operator: operator,
		Left:     left,
		Right:    right,
	}
}

func (n *BinOpNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitBinOp(n)
}

// UnaryOpNode - unary operation
type UnaryOpNode struct {
	BaseDecoratedNode
	Operator string
	Operand  DecoratedNode
}

func NewUnaryOpNode(operator string, operand DecoratedNode) *UnaryOpNode {
	return &UnaryOpNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Operator: operator,
		Operand:  operand,
	}
}

func (n *UnaryOpNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitUnaryOp(n)
}

// VarNode - variable reference
type VarNode struct {
	BaseDecoratedNode
	Name      string
	IsLValue  bool
	IsIndexed bool
	Index     DecoratedNode
}

func NewVarNode(name string) *VarNode {
	return &VarNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Name:     name,
		IsLValue: false,
	}
}

func (n *VarNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitVar(n)
}

// NumberNode - number literal
type NumberNode struct {
	BaseDecoratedNode
	Value int
}

func NewNumberNode(value int) *NumberNode {
	return &NumberNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeInteger,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Value: value,
	}
}

func (n *NumberNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitNumber(n)
}

// RealNode - real number literal
type RealNode struct {
	BaseDecoratedNode
	Value float64
}

func NewRealNode(value float64) *RealNode {
	return &RealNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeReal,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Value: value,
	}
}

func (n *RealNode) Accept(visitor DecoratedNodeVisitor) {
	// For now, treat like NumberNode
	visitor.VisitNumber(nil)
}

// StringNode - string literal
type StringNode struct {
	BaseDecoratedNode
	Value string
}

func NewStringNode(value string) *StringNode {
	return &StringNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeChar,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Value: value,
	}
}

func (n *StringNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitString(n)
}

// BooleanNode - boolean literal
type BooleanNode struct {
	BaseDecoratedNode
	Value bool
}

func NewBooleanNode(value bool) *BooleanNode {
	return &BooleanNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeBoolean,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Value: value,
	}
}

func (n *BooleanNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitBoolean(n)
}

// CharNode - character literal
type CharNode struct {
	BaseDecoratedNode
	Value rune
}

func NewCharNode(value rune) *CharNode {
	return &CharNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeChar,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Value: value,
	}
}

func (n *CharNode) Accept(visitor DecoratedNodeVisitor) {
	// For now, no visitor pattern implementation
}

// ProcCallNode - procedure call
type ProcCallNode struct {
	BaseDecoratedNode
	Name      string
	Arguments []DecoratedNode
	IsBuiltIn bool
}

func NewProcCallNode(name string, arguments []DecoratedNode) *ProcCallNode {
	return &ProcCallNode{
		BaseDecoratedNode: BaseDecoratedNode{
			TabIndex: -1,
			Type:     TypeNone,
			Ref:      -1,
			Errors:   make([]string, 0),
			Warnings: make([]string, 0),
		},
		Name:      name,
		Arguments: arguments,
		IsBuiltIn: false,
	}
}

func (n *ProcCallNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitProcCall(n)
}

// IfNode - if statement
type IfNode struct {
	BaseDecoratedNode
	Condition DecoratedNode
	ThenStmt  DecoratedNode
	ElseStmt  DecoratedNode
}

func (n *IfNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitIf(n)
}

// WhileNode - while loop
type WhileNode struct {
	BaseDecoratedNode
	Condition DecoratedNode
	Body      DecoratedNode
}

func (n *WhileNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitWhile(n)
}

// ForNode - for loop
type ForNode struct {
	BaseDecoratedNode
	Variable   DecoratedNode
	StartValue DecoratedNode
	EndValue   DecoratedNode
	Body       DecoratedNode
	IsDownTo   bool
}

func (n *ForNode) Accept(visitor DecoratedNodeVisitor) {
	visitor.VisitFor(n)
}

// ========== PRINT FUNCTION ==========

// PrintDecoratedAST prints the decorated AST in a tree format with visual connectors
func PrintDecoratedAST(node DecoratedNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	// Print current node connector
	connector := "+-- "
	if isLast {
		connector = "\\-- "
	}
	if prefix == "" {
		connector = ""
	}

	switch n := node.(type) {
	case *ProgramNode:
		fmt.Printf("%s%sProgramNode(name: '%s')\n", prefix, connector, n.Name)

		newPrefix := prefix
		if prefix != "" {
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "|   "
			}
		}

		// Print declarations and block with connectors
		hasBlock := n.Block != nil

		if n.Declarations != nil {
			fmt.Printf("%s|\n", newPrefix)
			fmt.Printf("%s+-- Declarations\n", newPrefix)
			declPrefix := newPrefix + "|   "
			printDeclarationList(n.Declarations, declPrefix)
		}

		if hasBlock {
			fmt.Printf("%s|\n", newPrefix)
			fmt.Printf("%s\\-- Block\n", newPrefix)
			blockPrefix := newPrefix + "    "
			printBlockStatements(n.Block, blockPrefix)
		}

	case *VarDeclNode:
		fmt.Printf("%s%sVarDecl(name: '%s', type: '%s')\n", prefix, connector, n.Name, n.Type)

	case *ConstDeclNode:
		fmt.Printf("%s%sConstDecl(name: '%s', value: %v, type: '%s')\n", prefix, connector, n.Name, n.Value, n.Type)

	case *SubprogramDeclNode:
		kind := "Procedure"
		if n.IsFunction {
			kind = "Function"
		}
		fmt.Printf("%s%s%sDecl(name: '%s', return_type: '%s')\n", prefix, connector, kind, n.Name, n.ReturnType)

	case *AssignNode:
		// Print assignment with multiline formatting for complex values
		targetStr := formatNodeInline(n.Target)

		// Check if we need multiline format
		if isBinOp(n.Value) {
			fmt.Printf("%s%sAssign(target: %s,\n", prefix, connector, targetStr)
			valueStr := formatBinOpMultiline(n.Value, prefix+getSpaces(len(connector))+"       ")
			fmt.Printf("%s%svalue: %s)\n", prefix, getSpaces(len(connector))+"       ", valueStr)
		} else {
			valueStr := formatNodeInline(n.Value)
			fmt.Printf("%s%sAssign(target: %s, value: %s)\n", prefix, connector, targetStr, valueStr)
		}

	case *BinOpNode:
		fmt.Printf("%s%sBinOp(op: '%s', left: %s, right: %s)\n",
			prefix, connector, n.Operator, formatNodeInline(n.Left), formatNodeInline(n.Right))

	case *UnaryOpNode:
		fmt.Printf("%s%sUnaryOp(op: '%s', operand: %s)\n", prefix, connector, n.Operator, formatNodeInline(n.Operand))

	case *VarNode:
		fmt.Printf("%s%sVar('%s')\n", prefix, connector, n.Name)

	case *NumberNode:
		fmt.Printf("%s%sNum(%d)\n", prefix, connector, n.Value)

	case *StringNode:
		fmt.Printf("%s%sString('%s')\n", prefix, connector, n.Value)

	case *CharNode:
		fmt.Printf("%s%sChar('%c')\n", prefix, connector, n.Value)

	case *BooleanNode:
		fmt.Printf("%s%sBool(%v)\n", prefix, connector, n.Value)

	case *ProcCallNode:
		// Format arguments inline
		args := make([]string, len(n.Arguments))
		for i, arg := range n.Arguments {
			args[i] = formatNodeInline(arg)
		}

		if len(args) == 0 {
			fmt.Printf("%s%sProcedureCall(name: '%s', args: [])\n", prefix, connector, n.Name)
		} else if len(args) == 1 {
			fmt.Printf("%s%sProcedureCall(name: '%s',\n", prefix, connector, n.Name)
			fmt.Printf("%s%sargs: [%s])\n", prefix, getSpaces(len(connector))+"              ", args[0])
		} else {
			// Multiline for multiple args
			fmt.Printf("%s%sProcedureCall(name: '%s',\n", prefix, connector, n.Name)
			argPrefix := prefix + getSpaces(len(connector)) + "              args: ["
			for i, arg := range args {
				if i == 0 {
					fmt.Printf("%s%s", argPrefix, arg)
				} else {
					fmt.Printf(", %s", arg)
				}
			}
			fmt.Printf("])\n")
		}

	case *IfNode:
		fmt.Printf("%s%sIf(condition: %s)\n", prefix, connector, formatNodeInline(n.Condition))

	case *WhileNode:
		fmt.Printf("%s%sWhile(condition: %s)\n", prefix, connector, formatNodeInline(n.Condition))

	case *ForNode:
		direction := "to"
		if n.IsDownTo {
			direction = "downto"
		}
		fmt.Printf("%s%sFor(var: %s, %s, start: %s, end: %s)\n",
			prefix, connector, formatNodeInline(n.Variable), direction,
			formatNodeInline(n.StartValue), formatNodeInline(n.EndValue))
	}
}

// Helper to print declaration list with proper connectors
func printDeclarationList(node DecoratedNode, prefix string) {
	if node == nil {
		return
	}

	if declList, ok := node.(*DeclarationListNode); ok {
		fmt.Printf("%s|\n", prefix)
		for i, decl := range declList.Declarations {
			isLast := i == len(declList.Declarations)-1
			PrintDecoratedAST(decl, prefix, isLast)
		}
	} else {
		PrintDecoratedAST(node, prefix, false)
	}
}

// Helper to print block statements with proper connectors
func printBlockStatements(node DecoratedNode, prefix string) {
	if node == nil {
		return
	}

	if block, ok := node.(*BlockNode); ok {
		fmt.Printf("%s|\n", prefix)
		for i, stmt := range block.Statements {
			isLast := i == len(block.Statements)-1
			PrintDecoratedAST(stmt, prefix, isLast)
			if !isLast {
				fmt.Printf("%s|\n", prefix)
			}
		}
	} else {
		PrintDecoratedAST(node, prefix, false)
	}
}

// Helper to format node inline
func formatNodeInline(node DecoratedNode) string {
	if node == nil {
		return "nil"
	}

	switch n := node.(type) {
	case *VarNode:
		return fmt.Sprintf("Var('%s')", n.Name)
	case *NumberNode:
		return fmt.Sprintf("Num(%d)", n.Value)
	case *StringNode:
		return fmt.Sprintf("String('%s')", n.Value)
	case *CharNode:
		return fmt.Sprintf("Char('%c')", n.Value)
	case *BooleanNode:
		return fmt.Sprintf("Bool(%v)", n.Value)
	case *BinOpNode:
		return fmt.Sprintf("BinOp(op: '%s', left: %s, right: %s)",
			n.Operator, formatNodeInline(n.Left), formatNodeInline(n.Right))
	case *UnaryOpNode:
		return fmt.Sprintf("UnaryOp(op: '%s', operand: %s)", n.Operator, formatNodeInline(n.Operand))
	case *ProcCallNode:
		// Format function/procedure call
		if len(n.Arguments) == 0 {
			return fmt.Sprintf("%s()", n.Name)
		}
		args := make([]string, len(n.Arguments))
		for i, arg := range n.Arguments {
			args[i] = formatNodeInline(arg)
		}
		argsStr := ""
		for i, arg := range args {
			if i > 0 {
				argsStr += ", "
			}
			argsStr += arg
		}
		return fmt.Sprintf("%s(%s)", n.Name, argsStr)
	default:
		return "<?>"
	}
}

// Helper to format BinOp with multiline
func formatBinOpMultiline(node DecoratedNode, prefix string) string {
	if binOp, ok := node.(*BinOpNode); ok {
		return fmt.Sprintf("BinOp(op: '%s',\n%sleft: %s,\n%sright: %s)",
			binOp.Operator, prefix, formatNodeInline(binOp.Left), prefix, formatNodeInline(binOp.Right))
	}
	return formatNodeInline(node)
}

// Helper to check if node is BinOp
func isBinOp(node DecoratedNode) bool {
	_, ok := node.(*BinOpNode)
	return ok
}

// Helper to generate spaces
func getSpaces(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += " "
	}
	return result
}
