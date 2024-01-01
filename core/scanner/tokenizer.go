package scanner

import (
	"SLang/model/token"
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

		switch ch {
		case '+':
			tokens = append(tokens, token.NewToken(token.T_PLUS))
			t.poll()
		case '-':
			tokens = append(tokens, token.NewToken(token.T_MINUS))
			t.poll()
		case '*':
			tokens = append(tokens, token.NewToken(token.T_STAR))
			t.poll()
		case '/':
			tokens = append(tokens, token.NewToken(token.T_SLASH))
			t.poll()
		default:
			if unicode.IsDigit(rune(ch)) {
				newToken := token.NewToken(token.T_INTLIT, t.scanNumber())
				tokens = append(tokens, newToken)
			} else {
				panic("Unrecognised character " + string(ch) + " on line " + strconv.Itoa(t.lineCount) + "\n")
			}
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
