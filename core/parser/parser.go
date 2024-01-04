package parser

import (
	"SLang/core/codegen"
	"SLang/core/scanner"
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

func (p *Parser) Parse(tokens []*token.Token) string {
	p.pos = 0
	if len(tokens) == 0 {
		return ""
	}
	return p.parseStatements(tokens)
}

func (p *Parser) parseStatements(tokens []*token.Token) string {
	generatedCode := ""
	for p.pos < len(tokens) {
		scanner.MatchToken(tokens[p.pos], token.T_PRINT)
		p.pos += 1

		if p.pos >= len(tokens) {
			panic("Expected expression")
		}

		treeRoot := p.parseBinaryExpr(0, tokens)
		res, retValue := treeRoot.CodeGen()
		generatedCode += res

		generatedCode += codegen.GenPrintInt(retValue)

		if p.pos >= len(tokens) {
			panic("Expected statement terminator")
		}

		scanner.MatchToken(tokens[p.pos], token.T_STATEMENT_TERMINATOR)
		p.pos += 1
	}
	return generatedCode
}

func (p *Parser) parseBinaryExpr(previousPrecedence int, tokens []*token.Token) astnode.ASTNode {
	var left, right astnode.ASTNode

	if tokens[p.pos].GetType() != token.T_INTLIT {
		panic("Invalid operand " + string(rune(tokens[p.pos].GetValue())))
	}

	// Process left which should always be an operand
	left = astnode.NewLeafNode(tokens[p.pos].GetValue())

	// Process next element which should always be an operator
	p.pos += 1
	if p.pos >= len(tokens) || tokens[p.pos].GetType() == token.T_STATEMENT_TERMINATOR {
		return left
	}

	operator := tokens[p.pos].GetType()

	for p.pos < len(tokens) && p.getOperatorPrecedence(operator) > previousPrecedence {
		// If current operator has higher precedence than previous operator, then process the right section
		// Else return the left
		p.pos += 1
		right = p.parseBinaryExpr(p.getOperatorPrecedence(operator), tokens)

		// Create new operator node and assign to left
		left = astnode.NewOpNode(operator, left, right)
		if p.pos >= len(tokens) || tokens[p.pos].GetType() == token.T_STATEMENT_TERMINATOR {
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
