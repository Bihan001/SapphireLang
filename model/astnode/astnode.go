package astnode

import (
	"SLang/core/codegen"
)

const (
	A_ADD = iota
	A_SUBTRACT
	A_MULTIPLY
	A_DIVIDE
	A_INTLIT
)

type ASTNode interface {
	CodeGen() (string, string)
}

type Const struct {
	ASTNode
	value int
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

func NewLeafNode(value int) ASTNode {
	return &Const{value: value}
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
	default:
		panic("Unknown operator")
	}
}

func (c *Const) CodeGen() (string, string) {
	return codegen.GetAllocateInstruction(c.value)
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
