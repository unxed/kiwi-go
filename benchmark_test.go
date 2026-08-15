package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

type view struct {
	left   *kiwi.Variable
	top    *kiwi.Variable
	width  *kiwi.Variable
	height *kiwi.Variable
	right  *kiwi.Variable
	bottom *kiwi.Variable
}

func newView() view {
	return view{
		left:   kiwi.NewVariable(),
		top:    kiwi.NewVariable(),
		width:  kiwi.NewVariable(),
		height: kiwi.NewVariable(),
		right:  kiwi.NewVariable(),
		bottom: kiwi.NewVariable(),
	}
}

type solverBenchData struct {
	solver    *kiwi.Solver
	superView view
	subView1  view
	subView2  view
}

func createKiwiSolver() solverBenchData {
	solver := kiwi.NewSolver()
	strength := kiwi.CreateStrength(0, 900, 1000)

	superView := newView()
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression(superView.left), kiwi.OpEq))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression(superView.top), kiwi.OpEq))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, superView.right}, superView.left, superView.width), kiwi.OpEq))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, superView.bottom}, superView.top, superView.height), kiwi.OpEq))

	solver.AddEditVariable(superView.width, kiwi.CreateStrength(999, 1000, 1000))
	solver.AddEditVariable(superView.height, kiwi.CreateStrength(999, 1000, 1000))
	solver.SuggestValue(superView.width, 300)
	solver.SuggestValue(superView.height, 200)

	subView1 := newView()
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.right}, subView1.left, subView1.width), kiwi.OpEq))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.bottom}, subView1.top, subView1.height), kiwi.OpEq))

	subView2 := newView()
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView2.right}, subView2.left, subView2.width), kiwi.OpEq))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView2.bottom}, subView2.top, subView2.height), kiwi.OpEq))

	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.left}, superView.left), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.top}, superView.top), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.bottom}, superView.bottom), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.width}, subView2.width), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView1.right}, subView2.left), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView2.right}, superView.right), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView2.top}, superView.top), kiwi.OpEq, nil, strength))
	solver.AddConstraint(kiwi.NewConstraint(kiwi.NewExpression([]any{-1, subView2.bottom}, superView.bottom), kiwi.OpEq, nil, strength))

	solver.UpdateVariables()

	return solverBenchData{
		solver:    solver,
		superView: superView,
		subView1:  subView1,
		subView2:  subView2,
	}
}

func TestBenchmarkSolverCorrectness(t *testing.T) {
	data := createKiwiSolver()
	if data.subView1.width.Value() != 150 {
		t.Errorf("expected subView1 width 150, got %f", data.subView1.width.Value())
	}
	if data.subView2.left.Value() != 150 {
		t.Errorf("expected subView2 left 150, got %f", data.subView2.left.Value())
	}
}

func BenchmarkCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = createKiwiSolver()
	}
}

func BenchmarkSolving(b *testing.B) {
	data := createKiwiSolver()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = data.solver.SuggestValue(data.superView.width, 100)
		_ = data.solver.SuggestValue(data.superView.height, 50)
		data.solver.UpdateVariables()
		_ = data.solver.SuggestValue(data.superView.width, 200)
		_ = data.solver.SuggestValue(data.superView.height, 500)
		data.solver.UpdateVariables()
	}
}
