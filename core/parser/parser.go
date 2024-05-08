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
	generatedCode := ""
	p.parseStatements(tokens, &generatedCode)
	return generatedCode
}

func (p *Parser) parseCompoundStatement(tokens []*token.Token, generatedCode *string) string {
	if p.pos >= len(tokens) {
		panic("Expected left brace")
	}

	p.lbrace(tokens)

	label := codegen.GetNewLabel()

	*generatedCode += label + ":\n"

	p.parseStatements(tokens, generatedCode)

	// p.rbrace(tokens)

	return label
}

func (p *Parser) parseStatements(tokens []*token.Token, generatedCode *string) {
	for p.pos < len(tokens) {

		registerService.FreeAllRegisters()

		switch tokens[p.pos].GetType() {
		case token.T_PRINT:
			p.parsePrintStatement(tokens, generatedCode)
			break
		case token.T_INT:
			p.parseVariableDeclaration(tokens, generatedCode)
			break
		case token.T_IDENTIFIER:
			p.parseAssignment(tokens, generatedCode)
			break
		case token.T_IF:
			p.parseIfStatement(tokens, generatedCode)
			break
		case token.T_RBRACE:
			p.rbrace(tokens)
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
}

func (p *Parser) parseIfStatement(tokens []*token.Token, generatedCode *string) {
	if p.pos >= len(tokens) {
		panic("Expected if statement")
	}
	scanner.MatchToken(tokens[p.pos], token.T_IF)
	p.pos += 1

	p.lparen(tokens)

	condition := p.parseBinaryExpr(0, tokens)
	res, retValue := condition.CodeGen()
	*generatedCode += res

	p.semi(tokens) // Ending binary expressions with ; for now
	p.rparen(tokens)

	returnLabel := codegen.GetNewLabel()
	
	ifLabel := p.parseCompoundStatement(tokens, generatedCode)
	*generatedCode += "\tjmp " + returnLabel + "\n"

	elseLabel := ""

	if p.pos < len(tokens) && tokens[p.pos].GetType() == token.T_ELSE {
		p.pos += 1
		elseLabel = p.parseCompoundStatement(tokens, generatedCode)
	} else {
		elseLabel = returnLabel
	}

	*generatedCode += codegen.GetIfInstruction(retValue, ifLabel, elseLabel)
	*generatedCode += returnLabel + ":\n"
}

func (p *Parser) parsePrintStatement(tokens []*token.Token, generatedCode *string) {
	if p.pos >= len(tokens) {
		panic("Expected print statement")
	}

	scanner.MatchToken(tokens[p.pos], token.T_PRINT)
	p.pos += 1

	treeRoot := p.parseBinaryExpr(0, tokens)
	res, retValue := treeRoot.CodeGen()
	*generatedCode += res

	*generatedCode += codegen.GenPrintInt(retValue)

	p.semi(tokens)
}

func (p *Parser) parseVariableDeclaration(tokens []*token.Token, generatedCode *string) string {
	if p.pos >= len(tokens) {
		panic("Expected variable declaration")
	}
	scanner.MatchToken(tokens[p.pos], token.T_INT)
	p.pos += 1

	// Get actual token value (string) from the tokenTable using token.GetValue()
	// Then create a new symbol in the symbol table using the fetched token value (string)
	identToken := tokens[p.pos]
	p.ident(tokens)
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

	if p.pos >= len(tokens) {
		panic("Expected statement terminator or assignment")
	}

	if tokens[p.pos].GetType() == token.T_ASSIGNMENT {
		p.pos -= 1
		p.parseAssignment(tokens, generatedCode)
	} else {
		p.semi(tokens)
	}

	return retValue
}

func (p *Parser) parseAssignment(tokens []*token.Token, generatedCode *string) string {
	if p.pos >= len(tokens) {
		panic("Expected identifier")
	}
	var left, right astnode.ASTNode

	identToken := tokens[p.pos]
	p.ident(tokens)
	tokenValue, ok := tokenTable.GetTokenAtPosition(identToken.GetValue())
	if !ok {
		panic("Unrecognised token")
	}
	identSymbolIdx := symbolTable.GetSymbol(tokenValue)
	if identSymbolIdx == -1 {
		panic("Undeclared variable " + tokenValue)
	}

	right = astnode.NewLeafNode(token.T_LIDENTIFIER, identSymbolIdx)

	p.assignment(tokens)

	left = p.parseBinaryExpr(0, tokens)

	tree := astnode.NewOpNode(token.T_ASSIGNMENT, right, left)
	resCode, retValue := tree.CodeGen()

	// generate code
	*generatedCode += resCode

	p.semi(tokens)

	return retValue
}

func (p *Parser) parseBinaryExpr(previousPrecedence int, tokens []*token.Token) astnode.ASTNode {
	if p.pos >= len(tokens) {
		panic("Expected expression")
	}

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

func (p *Parser) semi(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected statement terminator")
	}
	scanner.MatchToken(tokens[p.pos], token.T_STATEMENT_TERMINATOR)
	p.pos += 1
}

func (p *Parser) assignment(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected assignment")
	}
	scanner.MatchToken(tokens[p.pos], token.T_ASSIGNMENT)
	p.pos += 1
}

func (p *Parser) ident(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected identifier")
	}
	scanner.MatchToken(tokens[p.pos], token.T_IDENTIFIER)
	p.pos += 1
}

func (p *Parser) lbrace(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected left brace")
	}
	scanner.MatchToken(tokens[p.pos], token.T_LBRACE)
	p.pos += 1
}


func (p *Parser) rbrace(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected right brace")
	}
	scanner.MatchToken(tokens[p.pos], token.T_RBRACE)
	p.pos += 1
}

func (p *Parser) lparen(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected left parenthesis")
	}
	scanner.MatchToken(tokens[p.pos], token.T_LPAREN)
	p.pos += 1
}

func (p *Parser) rparen(tokens []*token.Token) {
	if p.pos >= len(tokens) {
		panic("Expected right parenthesis")
	}
	scanner.MatchToken(tokens[p.pos], token.T_RPAREN)
	p.pos += 1
}
