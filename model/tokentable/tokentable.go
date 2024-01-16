package tokentable

/*
TokenTable is similar to SymbolTable but is only used by the tokenizer.
It simply maps identifiers to integers since token.Token can hold only int values and not strings.
SymbolTable is filled later by the parser with actual logic (declara
*/
type TokenTable struct {
	tokens []string
}

const MAX_TOKENS_IN_TABLE = 1024

var tokenTable *TokenTable = nil

func NewTokenTable() *TokenTable {
	if tokenTable == nil {
		tokenTable = &TokenTable{tokens: make([]string, 0)}
	}
	return tokenTable
}

func (tt *TokenTable) AddToken(token string) int {
	if len(tt.tokens) >= MAX_TOKENS_IN_TABLE {
		panic("Token table full!")
	}

	for i, val := range tt.tokens {
		if val == token {
			return i
		}
	}

	tt.tokens = append(tt.tokens, token)

	return len(tt.tokens) - 1
}

func (tt *TokenTable) GetToken(token string) int {
	for i, val := range tt.tokens {
		if token == val {
			return i
		}
	}
	return -1
}

func (tt *TokenTable) GetTokenAtPosition(pos int) (string, bool) {
	if pos < 0 || pos >= len(tt.tokens) {
		return "", false
	}
	return tt.tokens[pos], true
}
