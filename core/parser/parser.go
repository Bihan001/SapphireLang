package parser

import (
	"SLang/core/codegen"
	"SLang/core/scanner"
	"SLang/model/astnode"
	"SLang/model/symboltable"
	"SLang/model/token"
	"SLang/model/tokentable"
	"SLang/util"
	"strconv"
)

// var operatorPrecedenceTable = []int{10, 10, 20, 20, 0}
var operatorPrecedenceTable = map[int]int{
	token.T_STAR:       10,
	token.T_SLASH:      10,
	token.T_PLUS:       9,
	token.T_MINUS:      9,
	token.T_LT:         8,
	token.T_LTE:        8,
	token.T_GT:         8,
	token.T_GTE:        8,
	token.T_EQ:         7,
	token.T_NE:         7,
	token.T_INTLIT:     0,
	token.T_IDENTIFIER: 0,
}

var registerService = util.GetNewRegisterService()

type Parser struct {
	pos int
}

var tokenTable *tokentable.TokenTable = tokentable.NewTokenTable()
var symbolTable *symboltable.SymbolTable = symboltable.NewSymbolTable()

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

		registerService.FreeAllRegisters()

		switch tokens[p.pos].GetType() {
		case token.T_PRINT:
			p.parsePrintStatement(tokens, &generatedCode)
			break
		case token.T_INT:
			p.parseVariableDeclaration(tokens, &generatedCode)
			break
		case token.T_IDENTIFIER:
			p.parseAssignment(tokens, &generatedCode)
			break
		default:
			{
				res, ok := token.TokenStringMap[tokens[p.pos].GetType()]
				if !ok {
					res = strconv.Itoa(tokens[p.pos].GetValue())
				}
				panic("Syntax Error at " + res)
			}
		}
	}
	return generatedCode
}

func (p *Parser) parsePrintStatement(tokens []*token.Token, generatedCode *string) {
	scanner.MatchToken(tokens[p.pos], token.T_PRINT)
	p.pos += 1

	if p.pos >= len(tokens) {
		panic("Expected expression")
	}

	treeRoot := p.parseBinaryExpr(0, tokens)
	res, retValue := treeRoot.CodeGen()
	*generatedCode += res

	*generatedCode += codegen.GenPrintInt(retValue)

	if p.pos >= len(tokens) {
		panic("Expected statement terminator")
	}

	p.semi(tokens[p.pos])
}

func (p *Parser) parseVariableDeclaration(tokens []*token.Token, generatedCode *string) string {
	scanner.MatchToken(tokens[p.pos], token.T_INT)
	p.pos += 1

	// Get actual token value (string) from the tokenTable using token.GetValue()
	// Then create a new symbol in the symbol table using the fetched token value (string)
	identToken := tokens[p.pos]
	p.ident(tokens[p.pos])
	tokenValue, ok := tokenTable.GetTokenAtPosition(identToken.GetValue())
	if !ok {
		panic("Unrecognised token " + tokenValue)
	}
	identifierIdx := symbolTable.AddSymbol(tokenValue)

	// generate code
	node := astnode.NewLeafNode(token.T_LIDENTIFIER, identifierIdx)
	resCode, retValue := node.CodeGen()
	*generatedCode += resCode

	symbolTable.UpdateStatus(tokenValue, symboltable.SYM_STATUS_DECLARED)

	if tokens[p.pos].GetType() == token.T_ASSIGNMENT {
		p.pos -= 1
		p.parseAssignment(tokens, generatedCode)
	} else {
		p.semi(tokens[p.pos])
	}

	return retValue
}

func (p *Parser) parseAssignment(tokens []*token.Token, generatedCode *string) string {
	var left, right astnode.ASTNode

	identToken := tokens[p.pos]
	p.ident(tokens[p.pos])
	tokenValue, ok := tokenTable.GetTokenAtPosition(identToken.GetValue())
	if !ok {
		panic("Unrecognised token")
	}
	identSymbolIdx := symbolTable.GetSymbol(tokenValue)
	if identSymbolIdx == -1 {
		panic("Undeclared variable " + tokenValue)
	}

	right = astnode.NewLeafNode(token.T_LIDENTIFIER, identSymbolIdx)

	p.assignment(tokens[p.pos])

	left = p.parseBinaryExpr(0, tokens)

	tree := astnode.NewOpNode(token.T_ASSIGNMENT, right, left)
	resCode, retValue := tree.CodeGen()

	// generate code
	*generatedCode += resCode

	p.semi(tokens[p.pos])

	return retValue
}

func (p *Parser) parseBinaryExpr(previousPrecedence int, tokens []*token.Token) astnode.ASTNode {
	var left, right astnode.ASTNode

	{
		switch tokens[p.pos].GetType() {
		case token.T_INTLIT:
			left = astnode.NewLeafNode(token.T_INTLIT, tokens[p.pos].GetValue())
			break
		case token.T_IDENTIFIER:
			tokenValue, ok := tokenTable.GetTokenAtPosition(tokens[p.pos].GetValue())
			if !ok {
				panic("Unrecognised token")
			}
			identSymbolIdx := symbolTable.GetSymbol(tokenValue)
			if identSymbolIdx == -1 {
				panic("Undeclared variable " + tokenValue)
			}
			left = astnode.NewLeafNode(token.T_IDENTIFIER, identSymbolIdx)
			break
		default:
			panic("invalid operand " + string(rune(tokens[p.pos].GetValue())))
		}
	}

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
	precedence, ok := operatorPrecedenceTable[operator]
	if !ok {
		panic("Syntax error - Invalid operator in expression")
	}
	return precedence
}

func (p *Parser) semi(actualToken *token.Token) {
	scanner.MatchToken(actualToken, token.T_STATEMENT_TERMINATOR)
	p.pos += 1
}

func (p *Parser) assignment(actualToken *token.Token) {
	scanner.MatchToken(actualToken, token.T_ASSIGNMENT)
	p.pos += 1
}

func (p *Parser) ident(actualToken *token.Token) {
	scanner.MatchToken(actualToken, token.T_IDENTIFIER)
	p.pos += 1
}
