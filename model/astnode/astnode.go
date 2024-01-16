package astnode

import (
	"SLang/core/codegen"
	"SLang/model/symboltable"
	"strconv"
)

const (
	A_ADD = iota
	A_SUBTRACT
	A_MULTIPLY
	A_DIVIDE
	A_INTLIT
	A_IDENT
	A_LIDENT
	A_ASSIGN
)

var symbolTable *symboltable.SymbolTable = symboltable.NewSymbolTable()

type ASTNode interface {
	CodeGen() (string, string)
}

type Const struct {
	ASTNode
	value int
}

type Ident struct {
	ASTNode
	symbolTableIdx int
}

type LIdent struct {
	ASTNode
	symbolTableIndex int
}

type Add struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type Sub struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type Mul struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type Div struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type Assign struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

func NewLeafNode(nodeType int, value int) ASTNode {
	switch nodeType {
	case A_INTLIT:
		return &Const{value: value}
	case A_IDENT:
		return &Ident{symbolTableIdx: value}
	case A_LIDENT:
		return &LIdent{symbolTableIndex: value}
	default:
		panic("Invalid node type")
	}
}

func NewOpNode(op int, left ASTNode, right ASTNode) ASTNode {
	switch op {
	case A_ADD:
		return &Add{left: left, right: right}
	case A_SUBTRACT:
		return &Sub{left: left, right: right}
	case A_MULTIPLY:
		return &Mul{left: left, right: right}
	case A_DIVIDE:
		return &Div{left: left, right: right}
	case A_ASSIGN:
		return &Assign{left: left, right: right}
	default:
		panic("Unknown operator")
	}
}

func (c *Const) CodeGen() (string, string) {
	return codegen.GetRegisterAllocateInstruction(strconv.Itoa(c.value))
}

func (l *LIdent) CodeGen() (string, string) {
	symbol, err := symbolTable.GetSymbolAtIndex(l.symbolTableIndex)

	if err != nil {
		panic(err)
	}

	if symbol.Status == symboltable.SYM_STATUS_PENDING {
		return codegen.GetVariableAllocateInstruction(symbol.Name)
	}

	return "", codegen.GetVariableInstructionFromSymbol(symbol.Name)
}

func (i *Ident) CodeGen() (string, string) {
	symbol, err := symbolTable.GetSymbolAtIndex(i.symbolTableIdx)

	if err != nil {
		panic(err)
	}

	return codegen.GetRegisterAllocateInstruction(codegen.GetVariableInstructionFromSymbol(symbol.Name))
}

func (a *Add) CodeGen() (string, string) {
	leftInstruction, leftValue := a.left.CodeGen()
	rightInstruction, rightValue := a.right.CodeGen()

	instruction := codegen.GetAddInstruction(leftValue, rightValue)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (s *Sub) CodeGen() (string, string) {
	leftInstruction, leftValue := s.left.CodeGen()
	rightInstruction, rightValue := s.right.CodeGen()

	instruction := codegen.GetSubtractInstruction(leftValue, rightValue)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (m *Mul) CodeGen() (string, string) {
	leftInstruction, leftValue := m.left.CodeGen()
	rightInstruction, rightValue := m.right.CodeGen()

	instruction := codegen.GetMultiplyInstruction(leftValue, rightValue)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (d *Div) CodeGen() (string, string) {
	leftInstruction, leftValue := d.left.CodeGen()
	rightInstruction, rightValue := d.right.CodeGen()

	instruction := codegen.GetDivideInstruction(leftValue, rightValue)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (a *Assign) CodeGen() (string, string) {
	leftInstruction, leftValue := a.left.CodeGen()
	rightInstruction, rightValue := a.right.CodeGen()

	instruction := codegen.GetVariableAssignInstruction(leftValue, rightValue)

	return leftInstruction + rightInstruction + instruction, leftValue
}
