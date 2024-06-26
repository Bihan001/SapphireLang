package scanner

import (
	"SLang/model/token"
	"SLang/model/tokentable"
	"strconv"
	"unicode"
)

type Tokenizer struct {
	pos         int
	lineCount   int
	input       *string
	inputLength int
}

var tokenTable *tokentable.TokenTable = tokentable.NewTokenTable()

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
		} else if ch == token.TokenStringMap[token.T_LT][0] {
			t.poll()
			if t.peek() == token.TokenStringMap[token.T_LTE][1] {
				tokens = append(tokens, token.NewToken(token.T_LTE))
				t.poll()
			} else {
				tokens = append(tokens, token.NewToken(token.T_LT))
			}
		} else if ch == token.TokenStringMap[token.T_GT][0] {
			t.poll()
			if t.peek() == token.TokenStringMap[token.T_GTE][1] {
				tokens = append(tokens, token.NewToken(token.T_GTE))
				t.poll()
			} else {
				tokens = append(tokens, token.NewToken(token.T_GT))
			}
		} else if ch == token.TokenStringMap[token.T_EQ][0] {
			t.poll()
			if t.peek() == token.TokenStringMap[token.T_EQ][1] {
				tokens = append(tokens, token.NewToken(token.T_EQ))
				t.poll()
			} else {
				tokens = append(tokens, token.NewToken(token.T_ASSIGNMENT))
			}
		} else if ch == token.TokenStringMap[token.T_NE][0] {
			t.poll()
			if t.peek() == token.TokenStringMap[token.T_NE][1] {
				tokens = append(tokens, token.NewToken(token.T_NE))
				t.poll()
			} else {
				panic("no symbol starting with " + string(ch) + " on line " + strconv.Itoa(t.lineCount) + " found\n")
			}
		} else if ch == token.TokenStringMap[token.T_LPAREN][0] {
			tokens = append(tokens, token.NewToken(token.T_LPAREN))
			t.poll()
		} else if ch == token.TokenStringMap[token.T_RPAREN][0] {
			tokens = append(tokens, token.NewToken(token.T_RPAREN))
			t.poll()
		} else if ch == token.TokenStringMap[token.T_LBRACE][0] {
			tokens = append(tokens, token.NewToken(token.T_LBRACE))
			t.poll()
		} else if ch == token.TokenStringMap[token.T_RBRACE][0] {
			tokens = append(tokens, token.NewToken(token.T_RBRACE))
			t.poll()
		} else if ch == token.TokenStringMap[token.T_STATEMENT_TERMINATOR][0] {
			tokens = append(tokens, token.NewToken(token.T_STATEMENT_TERMINATOR))
			t.poll()
		} else if unicode.IsDigit(rune(ch)) {
			newToken := token.NewToken(token.T_INTLIT, t.scanNumber())
			tokens = append(tokens, newToken)
		} else if unicode.IsLetter(rune(ch)) || ch == '_' {
			identifier := t.scanIdentifier()
			keyword := t.keyword(&identifier)

			if keyword != nil {
				tokens = append(tokens, keyword)
			} else {
				tokens = append(tokens, token.NewToken(token.T_IDENTIFIER, tokenTable.AddToken(identifier)))
			}
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
func (t *Tokenizer) keyword(s *string) *token.Token {
	if *s == token.TokenStringMap[token.T_PRINT] {
		return token.NewToken(token.T_PRINT)
	}
	if *s == token.TokenStringMap[token.T_INT] {
		return token.NewToken(token.T_INT)
	}
	if *s == token.TokenStringMap[token.T_IF] {
		return token.NewToken(token.T_IF)
	}
	if *s == token.TokenStringMap[token.T_ELSE] {
		return token.NewToken(token.T_ELSE)
	}
	return nil
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
