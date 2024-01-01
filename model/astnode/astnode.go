package astnode

import (
	"SLang/util"
	"fmt"
)

const (
	A_ADD = iota
	A_SUBTRACT
	A_MULTIPLY
	A_DIVIDE
	A_INTLIT
)

var irVariableService *util.IRVariableService = util.GetNewIRVariableService()

type ASTNode interface {
	//Eval() int
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

/*
func (c *Const) Eval() int {
	return c.value
}

func (a *Add) Eval() int {
	return a.left.Eval() + a.right.Eval()
}

func (s *Sub) Eval() int {
	return s.left.Eval() - s.right.Eval()
}

func (m *Mul) Eval() int {
	return m.left.Eval() * m.right.Eval()
}

func (d *Div) Eval() int {
	return d.left.Eval() / d.right.Eval()
}
*/

func (c *Const) CodeGen() (string, string) {
	allc := irVariableService.GetNewAllocation()
	str := fmt.Sprintf("%s = add i32 0, %d\n", allc, c.value)
	return str, allc
}

func (a *Add) CodeGen() (string, string) {
	lstr, lval := a.left.CodeGen()
	rstr, rval := a.right.CodeGen()
	allc := irVariableService.GetNewAllocation()
	str := fmt.Sprintf("%s = add i32 %s, %s\n", allc, lval, rval)
	err1 := irVariableService.FreeAllocation(lval)
	err2 := irVariableService.FreeAllocation(rval)
	if err1 != nil || err2 != nil {
		if err1 != nil {
			panic(err1)
		} else {
			panic(err2)
		}
	}
	return lstr + rstr + str, allc
}

func (s *Sub) CodeGen() (string, string) {
	lstr, lval := s.left.CodeGen()
	rstr, rval := s.right.CodeGen()
	allc := irVariableService.GetNewAllocation()
	str := fmt.Sprintf("%s = sub i32 %s, %s\n", allc, lval, rval)
	err1 := irVariableService.FreeAllocation(lval)
	err2 := irVariableService.FreeAllocation(rval)
	if err1 != nil || err2 != nil {
		if err1 != nil {
			panic(err1)
		} else {
			panic(err2)
		}
	}
	return lstr + rstr + str, allc
}

func (m *Mul) CodeGen() (string, string) {
	lstr, lval := m.left.CodeGen()
	rstr, rval := m.right.CodeGen()
	allc := irVariableService.GetNewAllocation()
	str := fmt.Sprintf("%s = mul i32 %s, %s\n", allc, lval, rval)
	err1 := irVariableService.FreeAllocation(lval)
	err2 := irVariableService.FreeAllocation(rval)
	if err1 != nil || err2 != nil {
		if err1 != nil {
			panic(err1)
		} else {
			panic(err2)
		}
	}
	return lstr + rstr + str, allc
}

func (d *Div) CodeGen() (string, string) {
	lstr, lval := d.left.CodeGen()
	rstr, rval := d.right.CodeGen()
	allc := irVariableService.GetNewAllocation()
	str := fmt.Sprintf("%s = sdiv i32 %s, %s\n", allc, lval, rval)
	err1 := irVariableService.FreeAllocation(lval)
	err2 := irVariableService.FreeAllocation(rval)
	if err1 != nil || err2 != nil {
		if err1 != nil {
			panic(err1)
		} else {
			panic(err2)
		}
	}
	return lstr + rstr + str, allc
}
