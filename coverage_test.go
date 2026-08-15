package kiwi_test

import (
	"encoding/json"
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestVariableJSON(t *testing.T) {
	v := kiwi.NewVariable("width")
	v.SetValue(100)

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("unexpected error marshaling variable: %v", err)
	}

	expected := `{"name":"width","value":100}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestNumericTypesInExpression(t *testing.T) {
	v := kiwi.NewVariable("x")

	e := kiwi.NewExpression(
		int8(1), int16(2), int32(3), int64(4),
		uint(5), uint8(6), uint16(7), uint32(8), uint64(9),
		float32(10.5), float64(11.5),
		[2]any{-1, v},
	)

	if e.Constant() != 67.0 {
		t.Errorf("expected constant 67, got %f", e.Constant())
	}
}

func TestOperatorAndSymbolString(t *testing.T) {
	opInvalid := kiwi.Operator(99)
	if opInvalid.String() != "Operator(99)" {
		t.Errorf("expected Operator(99), got %s", opInvalid.String())
	}

	sym := kiwi.NewSymbol(kiwi.SymbolSlack, 42)
	if sym.String() != "Symbol(42, 2)" {
		t.Errorf("expected Symbol(42, 2), got %s", sym.String())
	}
}

func TestSolverMaxIterationsLimit(t *testing.T) {
	solver := kiwi.NewSolver()
	solver.MaxIterations = 0 // force iteration limit error

	left := kiwi.NewVariable("left")
	width := kiwi.NewVariable("width")
	right := kiwi.NewVariable("right")

	_ = solver.AddEditVariable(left, kiwi.StrengthStrong)
	_ = solver.AddEditVariable(width, kiwi.StrengthStrong)

	cn := kiwi.NewConstraint(right, kiwi.OpEq, left.Plus(width))
	err := solver.AddConstraint(cn)
	if err == nil || err.Error() != "solver iterations exceeded" {
		t.Errorf("expected 'solver iterations exceeded' error, got %v", err)
	}
}
func TestSolverUnknownErrors(t *testing.T) {
	solver := kiwi.NewSolver()
	v := kiwi.NewVariable("v")
	cn := kiwi.NewConstraint(v, kiwi.OpEq)

	err := solver.RemoveConstraint(cn)
	if err == nil || err.Error() != "unknown constraint" {
		t.Errorf("expected 'unknown constraint', got %v", err)
	}

	err = solver.RemoveEditVariable(v)
	if err == nil || err.Error() != "unknown edit variable" {
		t.Errorf("expected 'unknown edit variable', got %v", err)
	}

	err = solver.SuggestValue(v, 100)
	if err == nil || err.Error() != "unknown edit variable" {
		t.Errorf("expected 'unknown edit variable', got %v", err)
	}
}

func TestIndexedMapCopyWithCustomVal(t *testing.T) {
	m := kiwi.NewIndexedMap[*dummyItem, string]()
	i := &dummyItem{id: 10}
	m.Insert(i, "hello")

	cp := m.Copy(func(val string) string {
		return val + "_copied"
	})

	p, ok := cp.Find(i)
	if !ok || p.Second != "hello_copied" {
		t.Errorf("expected 'hello_copied', got %v, ok=%v", p, ok)
	}
}
func TestRemoveNonBasicConstraint(t *testing.T) {
	solver := kiwi.NewSolver()
	x := kiwi.NewVariable("x")
	y := kiwi.NewVariable("y")

	_ = solver.AddEditVariable(x, kiwi.StrengthStrong)
	_ = solver.AddEditVariable(y, kiwi.StrengthStrong)
	_ = solver.SuggestValue(x, 100)
	_ = solver.SuggestValue(y, 200)

	cn := kiwi.NewConstraint(x, kiwi.OpLe, y, kiwi.StrengthMedium)
	err := solver.AddConstraint(cn)
	if err != nil {
		t.Fatalf("failed adding constraint: %v", err)
	}

	solver.UpdateVariables()

	err = solver.RemoveConstraint(cn)
	if err != nil {
		t.Fatalf("failed removing constraint: %v", err)
	}

	if solver.HasConstraint(cn) {
		t.Errorf("expected constraint to be removed")
	}
}

func TestInequalityOperators(t *testing.T) {
	solver := kiwi.NewSolver()
	x := kiwi.NewVariable("x")

	cnGe := kiwi.NewConstraint(x, kiwi.OpGe, 50, kiwi.StrengthMedium)
	err := solver.AddConstraint(cnGe)
	if err != nil {
		t.Fatalf("failed adding OpGe constraint: %v", err)
	}

	cnLe := kiwi.NewConstraint(x, kiwi.OpLe, 100, kiwi.StrengthMedium)
	err = solver.AddConstraint(cnLe)
	if err != nil {
		t.Fatalf("failed adding OpLe constraint: %v", err)
	}

	solver.UpdateVariables()
	if x.Value() < 50 || x.Value() > 100 {
		t.Errorf("expected x in [50, 100], got %f", x.Value())
	}
}
func TestMustCreateConstraint(t *testing.T) {
	solver := kiwi.NewSolver()
	x := kiwi.NewVariable("x")

	cn := solver.MustCreateConstraint(x, kiwi.OpEq, 10)
	if cn == nil || !solver.HasConstraint(cn) {
		t.Errorf("expected constraint created and added")
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic on duplicate MustCreateConstraint")
		}
	}()
	_ = solver.MustCreateConstraint(x, kiwi.OpEq, 10)
}

func TestRowRemoveSymbol(t *testing.T) {
	r := kiwi.NewRow(10)
	s := kiwi.NewSymbol(kiwi.SymbolSlack, 1)
	r.InsertSymbol(s, 5.0)

	if r.CoefficientFor(s) != 5.0 {
		t.Errorf("expected coefficient 5.0")
	}

	r.RemoveSymbol(s)
	if r.CoefficientFor(s) != 0.0 {
		t.Errorf("expected symbol removed")
	}
}
