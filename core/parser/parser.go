package parser

import (
	"SLang/model/astnode"
	"SLang/model/token"
)

var operatorPrecedenceTable = []int{10, 10, 20, 20, 0}

type Parser struct {
	pos int
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(tokens []*token.Token) astnode.ASTNode {
	p.pos = 0
	if len(tokens) == 0 {
		return nil
	}
	return p.parse(0, tokens)
}

func (p *Parser) parse(previousPrecedence int, tokens []*token.Token) astnode.ASTNode {
	var left, right astnode.ASTNode

	if tokens[p.pos].GetType() != token.T_INTLIT {
		panic("Invalid operand " + string(rune(tokens[p.pos].GetValue())))
	}

	// Process left which should always be an operand
	left = astnode.NewLeafNode(tokens[p.pos].GetValue())

	// Process next element which should always be an operator
	p.pos += 1
	if p.pos >= len(tokens) {
		return left
	}

	operator := tokens[p.pos].GetType()

	for p.pos < len(tokens) && p.getOperatorPrecedence(operator) > previousPrecedence {
		// If current operator has higher precedence than previous operator, then process the right section
		// Else return the left
		p.pos += 1
		right = p.parse(p.getOperatorPrecedence(operator), tokens)

		// Create new operator node and assign to left
		left = astnode.NewOpNode(operator, left, right)
		if p.pos >= len(tokens) {
			break
		}
		operator = tokens[p.pos].GetType()
	}

	return left
}

func (p *Parser) getOperatorPrecedence(operator int) int {
	if operator == token.T_INTLIT {
		panic("Syntax error")
	}
	return operatorPrecedenceTable[operator]
}
