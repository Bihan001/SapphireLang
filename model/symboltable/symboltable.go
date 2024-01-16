package symboltable

import (
	"errors"
	"fmt"
)

const (
	SYM_STATUS_PENDING = iota
	SYM_STATUS_DECLARED
)

type Symbol struct {
	Name   string
	Status int
}

type SymbolTable struct {
	symbols []*Symbol
}

const MAX_SYMBOLS_IN_TABLE = 1024

var symbolTable *SymbolTable = nil

func NewSymbolTable() *SymbolTable {
	if symbolTable == nil {
		symbolTable = &SymbolTable{symbols: make([]*Symbol, 0)}
	}
	return symbolTable
}

func (st *SymbolTable) AddSymbol(symbol string) int {
	if len(st.symbols) >= MAX_SYMBOLS_IN_TABLE {
		panic("Symbol table full!")
	}

	for _, val := range st.symbols {
		if val.Name == symbol {
			panic(fmt.Sprintf("%s symbol already declared", symbol))
		}
	}

	st.symbols = append(st.symbols, &Symbol{Name: symbol, Status: SYM_STATUS_PENDING})

	return len(st.symbols) - 1
}

func (st *SymbolTable) UpdateStatus(symbol string, status int) {
	i := -1

	for idx, s := range st.symbols {
		if s.Name == symbol {
			i = idx
			break
		}
	}

	if i == -1 {
		panic("No symbol found")
	}

	st.symbols[i].Status = status
}

func (st *SymbolTable) GetSymbol(symbol string) int {
	for i, val := range st.symbols {
		if symbol == val.Name {
			return i
		}
	}
	return -1
}

func (st *SymbolTable) GetSymbolAtIndex(index int) (*Symbol, error) {
	if index < 0 || index >= len(st.symbols) {
		return nil, errors.New("index out of bounds")
	}
	return st.symbols[index], nil
}
