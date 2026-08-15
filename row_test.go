package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestRowOperations(t *testing.T) {
	r := kiwi.NewRow(10.0)
	if r.Constant() != 10.0 {
		t.Errorf("expected constant 10.0, got %f", r.Constant())
	}

	s1 := kiwi.NewSymbol(kiwi.SymbolExternal, 1)
	s2 := kiwi.NewSymbol(kiwi.SymbolSlack, 2)
	sDummy := kiwi.NewSymbol(kiwi.SymbolDummy, 3)

	r.InsertSymbol(s1, 2.0)
	r.InsertSymbol(s2, -4.0)

	if r.CoefficientFor(s1) != 2.0 {
		t.Errorf("expected coefficient 2.0 for s1, got %f", r.CoefficientFor(s1))
	}
	if r.CoefficientFor(s2) != -4.0 {
		t.Errorf("expected coefficient -4.0 for s2, got %f", r.CoefficientFor(s2))
	}

	if r.AllDummies() {
		t.Errorf("expected AllDummies to be false")
	}

	r2 := kiwi.NewRow(0.0)
	r2.InsertSymbol(sDummy, 1.0)
	if !r2.AllDummies() {
		t.Errorf("expected AllDummies to be true for dummy row")
	}

	cp := r.Copy()
	if cp.Constant() != 10.0 || cp.CoefficientFor(s1) != 2.0 {
		t.Errorf("expected copied row to match original")
	}

	r.ReverseSign()
	if r.Constant() != -10.0 || r.CoefficientFor(s1) != -2.0 {
		t.Errorf("expected reversed sign")
	}

	// SolveFor
	r3 := kiwi.NewRow(10.0)
	r3.InsertSymbol(s1, 2.0)
	r3.SolveFor(s1)
	// 2*s1 + 10 = 0 => s1 = -5
	if r3.Constant() != -5.0 {
		t.Errorf("expected solved constant -5.0, got %f", r3.Constant())
	}

	// SolveForEx
	r4 := kiwi.NewRow(6.0)
	r4.InsertSymbol(s2, 2.0) // 2*s2 + 6 = 0
	r4.SolveForEx(s1, s2)
	// s1 = s2 / (-2) - 3...
	if r4.CoefficientFor(s1) == 0 {
		t.Errorf("expected s1 in cell map")
	}
}
