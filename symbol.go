package kiwi

import "fmt"

// SymbolType defines the available symbol types in the solver.
type SymbolType int

const (
	SymbolInvalid  SymbolType = 0
	SymbolExternal SymbolType = 1
	SymbolSlack    SymbolType = 2
	SymbolError    SymbolType = 3
	SymbolDummy    SymbolType = 4
)

// Symbol represents an internal symbol in the constraint solver.
type Symbol struct {
	id      int
	symType SymbolType
}

// NewSymbol constructs a new Symbol instance.
func NewSymbol(symType SymbolType, id int) *Symbol {
	return &Symbol{
		id:      id,
		symType: symType,
	}
}

// ID returns the unique id number of the symbol.
func (s *Symbol) ID() int {
	return s.id
}

// Type returns the type of the symbol.
func (s *Symbol) Type() SymbolType {
	return s.symType
}

func (s *Symbol) String() string {
	return fmt.Sprintf("Symbol(%d, %d)", s.id, s.symType)
}

var invalidSymbol = NewSymbol(SymbolInvalid, -1)
