package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func BenchmarkApportionSum10Cols(b *testing.B) {
	vars := make([]*kiwi.Variable, 10)
	floats := make(map[*kiwi.Variable]float64)
	for i := 0; i < 10; i++ {
		v := kiwi.NewVariable()
		vars[i] = v
		floats[v] = 12.3456
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kiwi.ApportionSum(vars, floats, 123)
	}
}

func BenchmarkDiscreteSolverTUILayout(b *testing.B) {
	screenW := kiwi.NewVariable("screen_w")
	sidebarW := kiwi.NewVariable("sidebar_w")
	col1 := kiwi.NewVariable("col1")
	col2 := kiwi.NewVariable("col2")
	col3 := kiwi.NewVariable("col3")

	ds := kiwi.NewDiscreteSolver()
	solver := ds.Solver()

	_ = solver.AddEditVariable(screenW, kiwi.StrengthStrong)
	_ = solver.SuggestValue(screenW, 120)

	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarW, kiwi.OpEq, screenW.Multiply(0.25)))
	_ = solver.AddConstraint(kiwi.NewConstraint(col1, kiwi.OpEq, col2))
	_ = solver.AddConstraint(kiwi.NewConstraint(col2, kiwi.OpEq, col3))
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarW.Plus(col1).Plus(col2).Plus(col3), kiwi.OpEq, screenW))

	ds.SetMinSize(sidebarW, 20)
	ds.AddApportionGroup(kiwi.ApportionGroup{
		Vars:      []*kiwi.Variable{col1, col2, col3},
		TargetVar: screenW,
	})
	ds.AddDirective(kiwi.SnapToGrid(sidebarW, 2))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = solver.SuggestValue(screenW, float64(100+(i%50)))
		_ = ds.SolveDiscrete()
	}
}
