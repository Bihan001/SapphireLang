package scanner

import (
	"SLang/model/token"
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type Tokenizer struct {
	pos         int
	lineCount   int
	input       *string
	inputLength int
}

func NewTokenizer(input *string) *Tokenizer {
	t := &Tokenizer{input: input, inputLength: len(*input), pos: 0, lineCount: 1}
	return t
}

func (t *Tokenizer) Tokenize() []*token.Token {
	var tokens []*token.Token

	for t.pos < t.inputLength {

		t.updateLines()
		t.skip()

		if t.isTokenized() {
			break
		}

		ch := t.peek()

		if ch == '+' {
			tokens = append(tokens, token.NewToken(token.T_PLUS))
			t.poll()
		} else if ch == '-' {
			tokens = append(tokens, token.NewToken(token.T_MINUS))
			t.poll()
		} else if ch == '*' {
			tokens = append(tokens, token.NewToken(token.T_STAR))
			t.poll()
		} else if ch == '/' {
			tokens = append(tokens, token.NewToken(token.T_SLASH))
			t.poll()
		} else if ch == token.TokenStringMap[token.T_STATEMENT_TERMINATOR][0] {
			tokens = append(tokens, token.NewToken(token.T_STATEMENT_TERMINATOR))
			t.poll()
		} else if unicode.IsDigit(rune(ch)) {
			newToken := token.NewToken(token.T_INTLIT, t.scanNumber())
			tokens = append(tokens, newToken)
		} else if unicode.IsLetter(rune(ch)) || ch == '_' {
			identifier := t.scanIdentifier()
			keyword, err := t.keyword(&identifier)
			if err != nil {
				panic(err)
			}
			tokens = append(tokens, keyword)
		} else {
			panic("Unrecognised character " + string(ch) + " on line " + strconv.Itoa(t.lineCount) + "\n")
		}
	}

	return tokens
}

func (t *Tokenizer) skip() {
	for t.pos < t.inputLength && (t.peek() == ' ' || t.peek() == '\t' || t.peek() == '\n' || t.peek() == '\r' || t.peek() == '\f') {
		t.poll()
	}
}

func (t *Tokenizer) peek() byte {
	if t.pos >= t.inputLength {
		panic("Index overflow while tokenizing")
	}
	return (*t.input)[t.pos]
}

func (t *Tokenizer) isTokenized() bool {
	return t.pos >= t.inputLength
}

func (t *Tokenizer) poll() {
	t.pos++
}

func (t *Tokenizer) updateLines() {
	if t.peek() == '\n' {
		t.lineCount += 1
	}
}

func (t *Tokenizer) scanNumber() int {
	n := 0

	for t.pos < t.inputLength && unicode.IsDigit(rune(t.peek())) {
		n = (n * 10) + (int(t.peek()) - '0')
		t.poll()
	}

	return n
}

// keyword This method checks whether the given string matches any existing identifier or not
func (t *Tokenizer) keyword(s *string) (*token.Token, error) {
	if *s == token.TokenStringMap[token.T_PRINT] {
		return token.NewToken(token.T_PRINT), nil
	}
	return nil, errors.New("Unrecognised symbol " + *s)
}

func (t *Tokenizer) scanIdentifier() string {
	i := 0
	identifier := ""

	for t.pos < t.inputLength && (unicode.IsLetter(rune(t.peek())) || unicode.IsDigit(rune(t.peek())) || t.peek() == '_') {
		if i >= token.MAX_TOKEN_LENGTH {
			panic("Token length greater than maximum token length: " + strconv.Itoa(token.MAX_TOKEN_LENGTH))
		}
		identifier += string(t.peek())
		i += 1
		t.poll()
	}

	return identifier
}

func MatchToken(actualToken *token.Token, expectedTokenType int) {
	if (*actualToken).GetType() == expectedTokenType {
		return
	}

	res, ok := token.TokenStringMap[actualToken.GetType()]
	if !ok {
		res = strconv.Itoa(actualToken.GetValue())
	}

	panic(fmt.Sprintf("Expected %s but got %s", token.TokenStringMap[expectedTokenType], res))
}
