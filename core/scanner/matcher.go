package scanner

import (
	"SLang/model/token"
	"fmt"
	"strconv"
)

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
