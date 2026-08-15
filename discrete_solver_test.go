package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestDiscreteSolverThreeColumns(t *testing.T) {
	c1 := kiwi.NewVariable("col1")
	c2 := kiwi.NewVariable("col2")
	c3 := kiwi.NewVariable("col3")
	total := kiwi.NewVariable("total")

	ds := kiwi.NewDiscreteSolver()
	solver := ds.Solver()

	_ = solver.AddEditVariable(total, kiwi.StrengthRequired)
	_ = solver.SuggestValue(total, 80)

	// col1 == col2, col2 == col3
	_ = solver.AddConstraint(kiwi.NewConstraint(c1, kiwi.OpEq, c2))
	_ = solver.AddConstraint(kiwi.NewConstraint(c2, kiwi.OpEq, c3))
	_ = solver.AddConstraint(kiwi.NewConstraint(c1.Plus(c2).Plus(c3), kiwi.OpEq, total))

	ds.AddApportionGroup(kiwi.ApportionGroup{
		Vars:      []*kiwi.Variable{c1, c2, c3},
		TargetVar: total,
	})

	res := ds.SolveDiscrete()

	if res.Get(total) != 80 {
		t.Errorf("expected total 80, got %d", res.Get(total))
	}
	if res.Sum(c1, c2, c3) != 80 {
		t.Errorf("expected sum 80, got %d", res.Sum(c1, c2, c3))
	}

	if res.Get(c1) < 26 || res.Get(c1) > 27 {
		t.Errorf("unexpected width for c1: %d", res.Get(c1))
	}
	if res.Get(c2) < 26 || res.Get(c2) > 27 {
		t.Errorf("unexpected width for c2: %d", res.Get(c2))
	}
	if res.Get(c3) < 26 || res.Get(c3) > 27 {
		t.Errorf("unexpected width for c3: %d", res.Get(c3))
	}
}

func TestDiscreteSolverTUIWindowLayout(t *testing.T) {
	screenWidth := kiwi.NewVariable("screen_width")
	sidebarWidth := kiwi.NewVariable("sidebar_width")
	mainWidth := kiwi.NewVariable("main_width")

	ds := kiwi.NewDiscreteSolver()
	solver := ds.Solver()

	_ = solver.AddEditVariable(screenWidth, kiwi.StrengthRequired)
	_ = solver.SuggestValue(screenWidth, 120)

	// sidebarWidth == 0.25 * screenWidth => 30
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarWidth, kiwi.OpEq, screenWidth.Multiply(0.25)))
	_ = solver.AddConstraint(kiwi.NewConstraint(sidebarWidth.Plus(mainWidth), kiwi.OpEq, screenWidth))

	ds.SetMinSize(sidebarWidth, 20)
	ds.SetMaxSize(sidebarWidth, 35)

	ds.AddApportionGroup(kiwi.ApportionGroup{
		Vars:      []*kiwi.Variable{sidebarWidth, mainWidth},
		TargetVar: screenWidth,
	})

	// Add a grid snapping directive to snap sidebar to even character count (e.g. CJK or double border)
	ds.AddDirective(kiwi.SnapToGrid(sidebarWidth, 2))

	res := ds.SolveDiscrete()

	if res.Get(screenWidth) != 120 {
		t.Errorf("expected screen width 120, got %d", res.Get(screenWidth))
	}
	if res.Sum(sidebarWidth, mainWidth) != 120 {
		t.Errorf("expected sidebar + main == 120, got %d", res.Sum(sidebarWidth, mainWidth))
	}
	if res.Get(sidebarWidth)%2 != 0 {
		t.Errorf("expected sidebar width to be snapped to even number, got %d", res.Get(sidebarWidth))
	}
}

func TestDiscreteSolverMinSizeClamp(t *testing.T) {
	ds := kiwi.NewDiscreteSolver()
	solver := ds.Solver()

	smallBox := kiwi.NewVariable("small")
	_ = solver.AddEditVariable(smallBox, kiwi.StrengthWeak)
	_ = solver.SuggestValue(smallBox, 0.4)

	ds.SetMinSize(smallBox, 3)

	res := ds.SolveDiscrete()

	if res.Get(smallBox) < 3 {
		t.Errorf("expected smallBox clamped to min 3, got %d", res.Get(smallBox))
	}
}
