package token

import "fmt"

const (
	T_PLUS = iota
	T_MINUS
	T_STAR
	T_SLASH
	T_INTLIT
	T_PRINT
	T_STATEMENT_TERMINATOR
	T_INT
	T_ASSIGNMENT
	T_IDENTIFIER
	T_LIDENTIFIER
	T_LT
	T_GT
	T_EQ
	T_NE
	T_LTE
	T_GTE
	T_IF
	T_ELSE
	T_LPAREN
	T_RPAREN
	T_LBRACE
	T_RBRACE
)

const MAX_TOKEN_LENGTH = 512

var TokenStringMap = map[int]string{
	T_PRINT:                "print",
	T_STATEMENT_TERMINATOR: ";",
	T_INT:                  "int",
	T_ASSIGNMENT:           "=",
	T_LT:                   "<",
	T_GT:                   ">",
	T_EQ:                   "==",
	T_NE:                   "!=",
	T_LTE:                  "<=",
	T_GTE:                  ">=",
	T_IF:                   "if",
	T_ELSE:                 "else",
	T_LPAREN:               "(",
	T_RPAREN:               ")",
	T_LBRACE:               "{",
	T_RBRACE:               "}",
}

type Token struct {
	token int
	value int
}

func NewToken(token int, value ...int) *Token {
	if len(value) != 0 {
		return &Token{token: token, value: value[0]}
	}
	return &Token{token: token}
}

func (token *Token) GetType() int {
	return token.token
}

func (token *Token) GetValue() int {
	return token.value
}

func (token *Token) Display() {
	fmt.Printf("{token: %d, intvalue: %d}\n", token.token, token.value)
}
