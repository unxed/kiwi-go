package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestSolverBasic(t *testing.T) {
	solver := kiwi.NewSolver()
	if solver == nil {
		t.Fatal("expected non-nil solver")
	}

	left := kiwi.NewVariable("left")
	width := kiwi.NewVariable("width")
	right := kiwi.NewVariable("right")

	err := solver.AddEditVariable(left, kiwi.StrengthStrong)
	if err != nil {
		t.Fatalf("unexpected error adding edit var: %v", err)
	}
	err = solver.AddEditVariable(width, kiwi.StrengthStrong)
	if err != nil {
		t.Fatalf("unexpected error adding edit var: %v", err)
	}

	err = solver.SuggestValue(left, 100)
	if err != nil {
		t.Fatalf("unexpected error suggesting value: %v", err)
	}
	err = solver.SuggestValue(width, 400)
	if err != nil {
		t.Fatalf("unexpected error suggesting value: %v", err)
	}

	// right == left + width => right - left - width == 0
	expr := kiwi.NewExpression([]any{-1, right}, left, width)
	cn := kiwi.NewConstraint(expr, kiwi.OpEq)
	err = solver.AddConstraint(cn)
	if err != nil {
		t.Fatalf("unexpected error adding constraint: %v", err)
	}

	solver.UpdateVariables()

	if right.Value() != 500 {
		t.Errorf("expected right value 500, got %f", right.Value())
	}

	if !solver.HasConstraint(cn) {
		t.Errorf("expected solver to have constraint")
	}
	if len(solver.GetConstraints()) == 0 {
		t.Errorf("expected non-empty constraints array")
	}
}
