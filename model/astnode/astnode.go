package astnode

import (
	"SLang/core/codegen"
	"SLang/model/symboltable"
	"SLang/model/token"
	"strconv"
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

type GreaterThan struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type LessThan struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type GreaterThanEquals struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type LessThanEquals struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type Equals struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

type NotEquals struct {
	ASTNode
	left  ASTNode
	right ASTNode
}

func NewLeafNode(nodeType int, value int) ASTNode {
	switch nodeType {
	case token.T_INTLIT:
		return &Const{value: value}
	case token.T_IDENTIFIER:
		return &Ident{symbolTableIdx: value}
	case token.T_LIDENTIFIER:
		return &LIdent{symbolTableIndex: value}
	default:
		panic("Invalid node type")
	}
}

func NewOpNode(op int, left ASTNode, right ASTNode) ASTNode {
	switch op {
	case token.T_PLUS:
		return &Add{left: left, right: right}
	case token.T_MINUS:
		return &Sub{left: left, right: right}
	case token.T_STAR:
		return &Mul{left: left, right: right}
	case token.T_SLASH:
		return &Div{left: left, right: right}
	case token.T_ASSIGNMENT:
		return &Assign{left: left, right: right}
	case token.T_GT:
		return &GreaterThan{left: left, right: right}
	case token.T_LT:
		return &LessThan{left: left, right: right}
	case token.T_GTE:
		return &GreaterThanEquals{left: left, right: right}
	case token.T_LTE:
		return &LessThanEquals{left: left, right: right}
	case token.T_EQ:
		return &Equals{left: left, right: right}
	case token.T_NE:
		return &NotEquals{left: left, right: right}
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

func (gt *GreaterThan) CodeGen() (string, string) {
	leftInstruction, leftValue := gt.left.CodeGen()
	rightInstruction, rightValue := gt.right.CodeGen()

	instruction := codegen.GetCompareInstruction(leftValue, rightValue, codegen.CMP_GT)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (gte *GreaterThanEquals) CodeGen() (string, string) {
	leftInstruction, leftValue := gte.left.CodeGen()
	rightInstruction, rightValue := gte.right.CodeGen()

	instruction := codegen.GetCompareInstruction(leftValue, rightValue, codegen.CMP_GTE)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (lt *LessThan) CodeGen() (string, string) {
	leftInstruction, leftValue := lt.left.CodeGen()
	rightInstruction, rightValue := lt.right.CodeGen()

	instruction := codegen.GetCompareInstruction(leftValue, rightValue, codegen.CMP_LT)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (lte *LessThanEquals) CodeGen() (string, string) {
	leftInstruction, leftValue := lte.left.CodeGen()
	rightInstruction, rightValue := lte.right.CodeGen()

	instruction := codegen.GetCompareInstruction(leftValue, rightValue, codegen.CMP_LTE)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (eq *Equals) CodeGen() (string, string) {
	leftInstruction, leftValue := eq.left.CodeGen()
	rightInstruction, rightValue := eq.right.CodeGen()

	instruction := codegen.GetCompareInstruction(leftValue, rightValue, codegen.CMP_EQ)

	return leftInstruction + rightInstruction + instruction, leftValue
}

func (ne *NotEquals) CodeGen() (string, string) {
	leftInstruction, leftValue := ne.left.CodeGen()
	rightInstruction, rightValue := ne.right.CodeGen()

	instruction := codegen.GetCompareInstruction(leftValue, rightValue, codegen.CMP_NE)

	return leftInstruction + rightInstruction + instruction, leftValue
}
