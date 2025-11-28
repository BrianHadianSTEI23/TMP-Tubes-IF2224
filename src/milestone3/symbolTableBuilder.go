package milestone3

import (
	"compiler/milestone2"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Builder untuk mem-populate symbol table dari AST
type SymbolTableBuilder struct {
	SymTable      *SymbolTable
	CurrentOffset int // Offset untuk variabel di stack frame
	Errors        []string
}

// Constructor
func NewSymbolTableBuilder() *SymbolTableBuilder {
	return &SymbolTableBuilder{
		SymTable:      NewSymbolTable(),
		CurrentOffset: 5, // header stack frame (g tau hrs ada atau g)
		Errors:        make([]string, 0),
	}
}

// Build symbol table dari AST
func (builder *SymbolTableBuilder) Build(ast *milestone2.AbstractSyntaxTree) error {
	if ast == nil {
		return fmt.Errorf("AST is nil")
	}

	// Traverse AST dan build symbol table
	err := builder.traverseProgram(ast)
	if err != nil {
		return err
	}

	// Return error jika ada semantic errors
	if len(builder.Errors) > 0 {
		return fmt.Errorf("symbol table build failed with %d error(s)", len(builder.Errors))
	}

	return nil
}

// Dapatkan symbol table yang sudah di-build
func (builder *SymbolTableBuilder) GetSymbolTable() *SymbolTable {
	return builder.SymTable
}

// Dapatkan daftar error
func (builder *SymbolTableBuilder) GetErrors() []string {
	return builder.Errors
}

// Traverse <program> node
func (builder *SymbolTableBuilder) traverseProgram(node *milestone2.AbstractSyntaxTree) error {
	if node.Value != "<program>" {
		return fmt.Errorf("expected <program> node, got %s", node.Value)
	}

	// Program structure: <program-header> <declaration-part> <compound-statement> DOT

	// Process program header (index 0) - extract program name
	if len(node.Children) > 0 {
		headerNode := node.Children[0]
		if headerNode.Value == "<program-header>" && len(headerNode.Children) > 1 {
			programNameNode := headerNode.Children[1]
			if strings.Contains(programNameNode.Value, "IDENTIFIER") {
				programName := extractValue(programNameNode.Value)
				// Add program name to symbol table
				builder.SymTable.Enter(
					programName,
					ObjProgram,
					TypeNone,
					0, // ref to block 0
					1,
					0,
				)
			}
		}
	}

	// Traverse declaration part (index 1)
	if len(node.Children) > 1 {
		err := builder.traverseDeclarationPart(node.Children[1])
		if err != nil {
			return err
		}
	}

	// Scan compound statement for procedure calls (index 2)
	if len(node.Children) > 2 && node.Children[2].Value == "<compound-statement>" {
		builder.scanProcedureCalls(node.Children[2])

		// Create block for main compound statement
		builder.SymTable.enterBlock()
		builder.SymTable.exitBlock()
	}

	return nil
}

// Traverse <declaration-part> node
func (builder *SymbolTableBuilder) traverseDeclarationPart(node *milestone2.AbstractSyntaxTree) error {
	if node.Value != "<declaration-part>" {
		return nil // Skip if not declaration part
	}

	// Traverse all children (const, type, var, subprogram declarations)
	for _, child := range node.Children {
		switch child.Value {
		case "<const-declaration>":
			err := builder.traverseConstDeclaration(child)
			if err != nil {
				return err
			}
		case "<type-declaration>":
			err := builder.traverseTypeDeclaration(child)
			if err != nil {
				return err
			}
		case "<var-declaration>":
			err := builder.traverseVarDeclaration(child)
			if err != nil {
				return err
			}
		case "<subprogram-declaration>":
			err := builder.traverseSubprogramDeclaration(child)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Traverse <const-declaration> node
func (builder *SymbolTableBuilder) traverseConstDeclaration(node *milestone2.AbstractSyntaxTree) error {
	// Skip "konstanta" keyword (index 0)
	// Traverse const definitions
	for i := 1; i < len(node.Children); i++ {
		child := node.Children[i]
		if child.Value == "<const-def>" {
			// Extract: IDENTIFIER = value ;
			if len(child.Children) >= 3 {
				identifierNode := child.Children[0]
				valueNode := child.Children[2]

				identifier := extractValue(identifierNode.Value)

				// Check duplicate
				if builder.SymTable.IsDeclaredInCurrentScope(identifier) {
					builder.Errors = append(builder.Errors,
						fmt.Sprintf("Duplicate constant declaration: %s", identifier))
					continue
				}

				// Handle STRING_LITERAL as array of char
				if strings.Contains(valueNode.Value, "STRING_LITERAL") {
					stringValue := extractValue(valueNode.Value)
					// Remove quotes
					if len(stringValue) >= 2 {
						stringValue = stringValue[1 : len(stringValue)-1]
					}
					stringLength := len(stringValue)

					// Create ATAB entry for string array
					atabIndex := builder.SymTable.EnterArray(
						0,             // xtyp (0 for base type)
						int(TypeChar), // etyp (element type)
						-1,            // eref (no further reference for char)
						1,             // low (string starts at index 1)
						stringLength,  // high (string length)
						1,             // elsz (char size = 1 byte)
					)

					// Enter constant as array type
					builder.SymTable.Enter(
						identifier,
						ObjConstant,
						TypeArray,
						atabIndex, // ref to ATAB
						1,         // normal
						0,         // address (not used for constants)
					)
				} else {
					value, typ := builder.extractConstValue(valueNode.Value)

					// Enter to symbol table
					builder.SymTable.Enter(
						identifier,
						ObjConstant,
						typ,
						-1,    // no ref
						1,     // normal
						value, // value stored in adr
					)
				}
			}
		}
	}

	return nil
}

// Traverse <type-declaration> node
func (builder *SymbolTableBuilder) traverseTypeDeclaration(node *milestone2.AbstractSyntaxTree) error {
	// Skip "tipe" keyword
	// Traverse type definitions: IDENTIFIER = <type> ;
	for i := 1; i < len(node.Children); i += 4 {
		if i+2 >= len(node.Children) {
			break
		}

		identifierNode := node.Children[i]
		typeNode := node.Children[i+2]

		identifier := extractValue(identifierNode.Value)

		// Check duplicate
		if builder.SymTable.IsDeclaredInCurrentScope(identifier) {
			builder.Errors = append(builder.Errors,
				fmt.Sprintf("Duplicate type declaration: %s", identifier))
			continue
		}

		// Process type
		typ, ref := builder.processType(typeNode)

		// Enter to symbol table
		builder.SymTable.Enter(
			identifier,
			ObjType,
			typ,
			ref,
			1,
			0, // types don't have address
		)
	}

	return nil
}

// Traverse <var-declaration> node
func (builder *SymbolTableBuilder) traverseVarDeclaration(node *milestone2.AbstractSyntaxTree) error {
	// Skip "variabel" keyword (index 0)
	// Traverse: <identifier-list> : <type> ;
	for i := 1; i < len(node.Children); i += 4 {
		if i+2 >= len(node.Children) {
			break
		}

		identifierListNode := node.Children[i]
		typeNode := node.Children[i+2]

		// Extract identifiers
		identifiers := builder.extractIdentifierList(identifierListNode)

		// Process type
		typ, ref := builder.processType(typeNode)
		typeSize := builder.SymTable.getTypeSize(typ, ref)

		// Enter each identifier
		for _, identifier := range identifiers {
			// Check duplicate
			if builder.SymTable.IsDeclaredInCurrentScope(identifier) {
				builder.Errors = append(builder.Errors,
					fmt.Sprintf("deklarasi variabel duplikat: %s", identifier))
				continue
			}

			// Enter to symbol table
			builder.SymTable.Enter(
				identifier,
				ObjVariable,
				typ,
				ref,
				1, // normal variable
				builder.CurrentOffset,
			)
			builder.SymTable.AddVariableSize(typeSize)
			builder.CurrentOffset += typeSize
		}
	}

	return nil
}

// Traverse <subprogram-declaration> node
func (builder *SymbolTableBuilder) traverseSubprogramDeclaration(node *milestone2.AbstractSyntaxTree) error {
	// Structure: (fungsi|prosedur) ID ( params ) (: type)? ; <declaration-part> <compound-statement> ;

	keywordNode := node.Children[0]
	nameNode := node.Children[1]

	isFungsi := strings.Contains(keywordNode.Value, "fungsi")
	name := extractValue(nameNode.Value)

	// Check duplicate
	// if builder.SymTable.IsDeclaredInCurrentScope(name) {
	// 	builder.Errors = append(builder.Errors,
	// 		fmt.Sprintf("duplikasi deklarasi subprogram: %s", name))
	// 	return nil
	// }

	// Create new block for subprogram
	blockIndex := builder.SymTable.EnterBlock()
	builder.SymTable.enterLevel()

	oldOffset := builder.CurrentOffset
	builder.CurrentOffset = 5 // reset ke base offset stack frame

	lastParamIndex := 0

	// Find parameter list (between LPARENTHESIS and RPARENTHESIS)
	paramListIndex := -1
	for i, child := range node.Children {
		if child.Value == "<parameter-list>" {
			paramListIndex = i
			break
		}
	}

	if paramListIndex > 0 {
		paramListNode := node.Children[paramListIndex]
		params := builder.extractParameters(paramListNode)

		for _, param := range params {
			size := builder.SymTable.getTypeSize(param.Type, param.Ref)
			lastParamIndex = builder.SymTable.Enter(
				param.Name,
				ObjVariable,
				param.Type,
				param.Ref,
				param.Nrm,
				builder.CurrentOffset,
			)
			builder.SymTable.AddParameterSize(size)
			builder.CurrentOffset += size
		}
	}

	// Update block parameter info
	builder.SymTable.UpdateBlockLastParam(blockIndex, lastParamIndex)

	// Determine return type for function
	returnType := TypeNone
	if isFungsi {
		// Find return type (after COLON)
		for i := 0; i < len(node.Children)-1; i++ {
			if strings.Contains(node.Children[i].Value, "COLON") && i+1 < len(node.Children) {
				if node.Children[i+1].Value == "<type>" {
					returnType, _ = builder.processType(node.Children[i+1])
					break
				}
			}
		}
	}

	// Enter subprogram to parent scope
	// objClass := ObjProcedure
	// if isFungsi {
	// 	objClass = ObjFunction
	// }

	// Temporarily exit to parent scope to enter the subprogram
	builder.SymTable.exitLevel()

	// subprogIndex := builder.SymTable.Enter(
	// 	name,
	// 	objClass,
	// 	returnType,
	// 	blockIndex,
	// 	1,
	// 	0, // address will be set during code generation
	// )

	builder.SymTable.enterLevel()

	// For functions, add implicit return variable (same name as function)
	if isFungsi {
		builder.SymTable.Enter(
			name,
			ObjVariable,
			returnType,
			-1,
			1,
			0,
		)
	}

	// Process local declarations
	// Find <declaration-part> in subprogram body
	for _, child := range node.Children {
		if child.Value == "<declaration-part>" {
			builder.traverseDeclarationPart(child)
			break
		}
	}

	// Exit subprogram scope
	builder.SymTable.exitLevel()
	builder.SymTable.exitBlock()
	builder.CurrentOffset = oldOffset

	return nil
}

// Process <type> node dan return TypeKind + ref
func (builder *SymbolTableBuilder) processType(node *milestone2.AbstractSyntaxTree) (TypeKind, int) {
	// Handle <array-type> directly (from type declarations)
	if node.Value == "<array-type>" {
		return builder.processArrayType(node)
	}

	// Handle <record-type> directly (from type declarations)
	if node.Value == "<record-type>" {
		return builder.processRecordType(node)
	}

	if node.Value != "<type>" {
		return TypeNone, -1
	}

	if len(node.Children) == 0 {
		return TypeNone, -1
	}

	child := node.Children[0]

	// Check if it's a built-in type
	if strings.Contains(child.Value, "KEYWORD") {
		typeStr := extractValue(child.Value)
		switch typeStr {
		case "integer":
			return TypeInteger, -1
		case "boolean":
			return TypeBoolean, -1
		case "char":
			return TypeChar, -1
			// case "real":
			// 	return TypeReal, -1
		}
	}

	// Check if it's array type
	if child.Value == "<array-type>" {
		return builder.processArrayType(child)
	}

	// Check if it's record type
	if child.Value == "<record-type>" {
		return builder.processRecordType(child)
	}

	// Check if it's user-defined type (IDENTIFIER)
	if strings.Contains(child.Value, "IDENTIFIER") {
		typeName := extractValue(child.Value)
		// Lookup type in symbol table
		idx, found := builder.SymTable.Lookup(typeName)
		if found && builder.SymTable.Tab[idx].Obj == ObjType {
			return builder.SymTable.Tab[idx].Type, builder.SymTable.Tab[idx].Ref
		}
		// Type not found - error
		builder.Errors = append(builder.Errors,
			fmt.Sprintf("Undefined type: %s", typeName))
		return TypeNone, -1
	}

	return TypeNone, -1
}

// Process <array-type> node
func (builder *SymbolTableBuilder) processArrayType(node *milestone2.AbstractSyntaxTree) (TypeKind, int) {
	// Structure: larik [ NUMBER .. NUMBER ] dari <type>

	// Extract bounds
	low := 0
	high := 0
	elementTypeNode := (*milestone2.AbstractSyntaxTree)(nil)

	for i, child := range node.Children {
		if strings.Contains(child.Value, "NUMBER") {
			val, _ := strconv.Atoi(extractValue(child.Value))
			if low == 0 {
				low = val
			} else {
				high = val
			}
		}
		if child.Value == "<type>" || child.Value == "<array-type>" {
			elementTypeNode = child
		}
		// Handle nested array
		if i < len(node.Children)-1 && strings.Contains(child.Value, "dari") {
			if node.Children[i+1].Value == "<array-type>" || node.Children[i+1].Value == "<type>" {
				elementTypeNode = node.Children[i+1]
			}
		}
	}

	// Process element type
	elemType, elemRef := TypeInteger, -1
	if elementTypeNode != nil {
		if elementTypeNode.Value == "<array-type>" {
			elemType, elemRef = builder.processArrayType(elementTypeNode)
		} else {
			elemType, elemRef = builder.processType(elementTypeNode)
		}
	}

	// Calculate element size
	elemSize := builder.SymTable.getTypeSize(elemType, elemRef)

	// Enter to array table
	atabIndex := builder.SymTable.EnterArray(
		int(TypeInteger), // index type (always integer)
		int(elemType),
		elemRef,
		low,
		high,
		elemSize,
	)

	return TypeArray, atabIndex
}

// Process <record-type> node
func (builder *SymbolTableBuilder) processRecordType(node *milestone2.AbstractSyntaxTree) (TypeKind, int) {
	// Structure: rekaman <field-list> selesai
	// Record type creates a block in btab (similar to procedure/function)

	// Save parent block
	parentBlock := builder.SymTable.CurrentBlock

	// Create new block for record
	blockIndex := builder.SymTable.EnterBlock()
	builder.SymTable.enterLevel()

	// Save old offset
	oldOffset := builder.CurrentOffset
	builder.CurrentOffset = 0

	// Process field declarations
	// Find <field-list> node (should be between 'rekaman' keyword and 'selesai')
	for _, child := range node.Children {
		if child.Value == "<field-list>" {
			err := builder.traverseFieldList(child)
			if err != nil {
				builder.Errors = append(builder.Errors,
					fmt.Sprintf("Error processing record fields: %v", err))
			}
			break
		}
	}

	// Update block size (vsze = total size of all fields)
	recordSize := builder.CurrentOffset
	if blockIndex >= 0 && blockIndex < len(builder.SymTable.Btab) {
		builder.SymTable.Btab[blockIndex].Vsze = recordSize
		// For records, lpar should be 0 (no parameters)
		builder.SymTable.Btab[blockIndex].Lpar = 0
	}

	// Exit record scope and restore parent block
	builder.SymTable.exitLevel()
	builder.SymTable.CurrentBlock = parentBlock
	builder.CurrentOffset = oldOffset

	return TypeRecord, blockIndex
}

// Traverse <field-list> node for record fields
func (builder *SymbolTableBuilder) traverseFieldList(node *milestone2.AbstractSyntaxTree) error {
	// Field list structure: <identifier-list> : <type> (; <identifier-list> : <type>)*
	// Similar to variable declarations

	for i := 0; i < len(node.Children); i++ {
		child := node.Children[i]

		// Skip semicolons
		if strings.Contains(child.Value, "SEMICOLON") {
			continue
		}

		// Process field declaration: <identifier-list> : <type>
		if child.Value == "<identifier-list>" {
			if i+2 >= len(node.Children) {
				continue // Not enough children for complete field declaration
			}

			identifierListNode := child
			// Skip colon (i+1)
			typeNode := node.Children[i+2]

			// Extract identifiers
			identifiers := builder.extractIdentifierList(identifierListNode)

			// Process type
			typ, ref := builder.processType(typeNode)
			typeSize := builder.SymTable.getTypeSize(typ, ref)

			// Enter each field
			for _, identifier := range identifiers {
				// Check duplicate
				if builder.SymTable.IsDeclaredInCurrentScope(identifier) {
					builder.Errors = append(builder.Errors,
						fmt.Sprintf("Duplicate field declaration: %s", identifier))
					continue
				}

				// Enter to symbol table as field
				builder.SymTable.Enter(
					identifier,
					ObjField,
					typ,
					ref,
					1, // normal field
					builder.CurrentOffset,
				)
				builder.SymTable.AddVariableSize(typeSize)
				builder.CurrentOffset += typeSize
			}

			// Skip to next field declaration (skip colon and type)
			i += 2
		}
	}

	return nil
}

// Scan AST untuk menemukan procedure calls dan tambahkan ke symbol table
func (builder *SymbolTableBuilder) scanProcedureCalls(node *milestone2.AbstractSyntaxTree) {
	if node == nil {
		return
	}

	// Check if this is a procedure call
	if node.Value == "<procedure-call>" && len(node.Children) > 0 {
		firstChild := node.Children[0]
		if strings.Contains(firstChild.Value, "IDENTIFIER") {
			procName := extractValue(firstChild.Value)

			// Check if it's a built-in procedure
			builtIns := map[string]bool{
				"writeln": true,
				"write":   true,
				"read":    true,
				"readln":  true,
			}

			if builtIns[procName] {
				// Check if already in symbol table
				if _, found := builder.SymTable.Lookup(procName); !found {
					// Add to symbol table
					builder.SymTable.Enter(
						procName,
						ObjProcedure,
						TypeNone,
						-1, // No block (built-in)
						1,
						-1, // Built-in, no address
					)
				}
			}
		}
	}

	for _, child := range node.Children {
		builder.scanProcedureCalls(child)
	}
}

// Extract list of identifiers dari <identifier-list> node
func (builder *SymbolTableBuilder) extractIdentifierList(node *milestone2.AbstractSyntaxTree) []string {
	identifiers := make([]string, 0)

	for _, child := range node.Children {
		if strings.Contains(child.Value, "IDENTIFIER") {
			identifiers = append(identifiers, extractValue(child.Value))
		}
	}

	return identifiers
}

// Struct untuk parameter
type Parameter struct {
	Name string
	Type TypeKind
	Ref  int
	Nrm  int
}

// Extract parameters dari <parameter-list> node
func (builder *SymbolTableBuilder) extractParameters(node *milestone2.AbstractSyntaxTree) []Parameter {
	params := make([]Parameter, 0)

	// Parameter structure: <identifier-list> : <type> (; <identifier-list> : <type>)*
	i := 0
	for i < len(node.Children) {
		isRef := false
		if strings.Contains(strings.ToLower(node.Children[i].Value), "var") {
			isRef = true
			i++
		}
		if i >= len(node.Children) {
			break
		}
		if node.Children[i].Value == "<identifier-list>" {
			idList := builder.extractIdentifierList(node.Children[i])

			if i+2 < len(node.Children) && node.Children[i+2].Value == "<type>" {
				typ, ref := builder.processType(node.Children[i+2])
				nrmVal := 1
				if isRef {
					nrmVal = 0
				}

				for _, id := range idList {
					params = append(params, Parameter{
						Name: id,
						Type: typ,
						Ref:  ref,
						Nrm:  nrmVal,
					})
				}
				i += 3
			} else {
				i++
			}
		} else {
			i++
		}
	}

	return params
}

// Extract nilai dan tipe dari constant value
func (builder *SymbolTableBuilder) extractConstValue(value string) (int, TypeKind) {
	if strings.Contains(value, "NUMBER") {
		val, _ := strconv.Atoi(extractValue(value))
		return val, TypeInteger
	}
	if strings.Contains(value, "STRING_LITERAL") {
		return 0, TypeChar // String treated as char array
	}
	if strings.Contains(value, "true") {
		return 1, TypeBoolean
	}
	if strings.Contains(value, "false") {
		return 0, TypeBoolean
	}
	return 0, TypeNone
}

// Extract value dari token string (e.g., "IDENTIFIER(x)" -> "x")
func extractValue(token string) string {
	re := regexp.MustCompile(`^[A-Z_]+\((.*)\)$`)
	matches := re.FindStringSubmatch(token)
	if len(matches) >= 2 {
		return matches[1]
	}
	return token
}
