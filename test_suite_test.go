package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestFullTestSuite(t *testing.T) {
	t.Run("kiwi create Solver", func(t *testing.T) {
		solver := kiwi.NewSolver()
		if solver == nil {
			t.Fatal("expected solver")
		}
	})

	t.Run("Variable test suite", func(t *testing.T) {
		solver := kiwi.NewSolver()
		v := kiwi.NewVariable()
		if v.Value() != 0 {
			t.Errorf("expected 0, got %f", v.Value())
		}

		v2 := kiwi.NewVariable("somename")
		if v2.Name() != "somename" {
			t.Errorf("expected 'somename', got %s", v2.Name())
		}

		v3 := kiwi.NewVariable()
		v3.SetName("skiwi")
		if v3.Name() != "skiwi" {
			t.Errorf("expected 'skiwi', got %s", v3.Name())
		}

		err := solver.AddEditVariable(v, kiwi.StrengthStrong)
		if err != nil || !solver.HasEditVariable(v) {
			t.Errorf("failed adding edit var")
		}

		err = solver.SuggestValue(v, 200)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		solver.UpdateVariables()
		if v.Value() != 200 {
			t.Errorf("expected 200, got %f", v.Value())
		}

		// Subscribe test
		subVar := kiwi.NewVariable()
		var val, prevVal float64
		callCount := 0
		subVar.Subscribe(func(value, previousValue float64) {
			callCount++
			val = value
			prevVal = previousValue
		})

		solver.AddEditVariable(subVar, kiwi.StrengthStrong)
		solver.SuggestValue(subVar, 400)
		solver.UpdateVariables()
		if val != 400 || prevVal != 0 {
			t.Errorf("expected 400, 0, got %f, %f", val, prevVal)
		}

		solver.SuggestValue(subVar, 500)
		solver.UpdateVariables()
		if val != 500 || prevVal != 400 {
			t.Errorf("expected 500, 400, got %f, %f", val, prevVal)
		}

		// Unsubscribe
		subVar.Unsubscribe()
		val = 0
		prevVal = 0
		solver.SuggestValue(subVar, 300)
		solver.UpdateVariables()
		if val != 0 || prevVal != 0 {
			t.Errorf("unsubscribe failed, callback was called")
		}

		// Remove edit var
		solver.RemoveEditVariable(v)
		if solver.HasEditVariable(v) {
			t.Errorf("expected edit var removed")
		}

		// Duplicate edit var error
		var4 := kiwi.NewVariable()
		solver.AddEditVariable(var4, kiwi.StrengthStrong)
		err = solver.AddEditVariable(var4, kiwi.StrengthStrong)
		if err == nil || err.Error() != "duplicate edit variable" {
			t.Errorf("expected 'duplicate edit variable' error, got %v", err)
		}

		// Unknown edit var error
		solver.RemoveEditVariable(var4)
		err = solver.RemoveEditVariable(var4)
		if err == nil || err.Error() != "unknown edit variable" {
			t.Errorf("expected 'unknown edit variable' error, got %v", err)
		}
	})

	t.Run("Constraint test suite", func(t *testing.T) {
		solver := kiwi.NewSolver()

		expr := kiwi.NewExpression(10)
		cn := kiwi.NewConstraint(expr, kiwi.OpEq)
		if cn.Expression() != expr {
			t.Errorf("expected matching expression")
		}

		cnGe := kiwi.NewConstraint(kiwi.NewExpression(10), kiwi.OpGe)
		if cnGe.Op() != kiwi.OpGe {
			t.Errorf("expected OpGe")
		}

		cnMed := kiwi.NewConstraint(kiwi.NewExpression(10), kiwi.OpLe, nil, kiwi.StrengthMedium)
		if cnMed.Strength() != kiwi.StrengthMedium {
			t.Errorf("expected StrengthMedium")
		}

		cnOpt := kiwi.NewConstraint(kiwi.NewExpression(1), kiwi.OpEq)
		if cnOpt.Strength() != kiwi.StrengthRequired {
			t.Errorf("expected StrengthRequired")
		}

		cn1 := kiwi.NewConstraint(kiwi.NewExpression(1, -1), kiwi.OpEq)
		if solver.HasConstraint(cn1) {
			t.Errorf("should not have cn1 yet")
		}
		solver.AddConstraint(cn1)
		if !solver.HasConstraint(cn1) {
			t.Errorf("should have cn1")
		}
		solver.RemoveConstraint(cn1)
		if solver.HasConstraint(cn1) {
			t.Errorf("should not have cn1")
		}

		// Duplicate constraint error
		solver.AddConstraint(cn1)
		err := solver.AddConstraint(cn1)
		if err == nil || err.Error() != "duplicate constraint" {
			t.Errorf("expected 'duplicate constraint', got %v", err)
		}

		// Unsatisfiable constraint (all numbers non-zero)
		cnUnsat := kiwi.NewConstraint(kiwi.NewExpression(1, -1, 10), kiwi.OpEq)
		err = solver.AddConstraint(cnUnsat)
		if err == nil || err.Error() != "unsatisfiable constraint" {
			t.Errorf("expected 'unsatisfiable constraint', got %v", err)
		}

		// Unsatisfiable constraint (conflicting variables)
		solver2 := kiwi.NewSolver()
		w1 := kiwi.NewVariable()
		w2 := kiwi.NewVariable()
		cnA := kiwi.NewConstraint(kiwi.NewExpression(w1, 100), kiwi.OpEq)
		solver2.AddConstraint(cnA)
		cnB := kiwi.NewConstraint(kiwi.NewExpression(w2, 100), kiwi.OpEq)
		solver2.AddConstraint(cnB)
		cnC := kiwi.NewConstraint(kiwi.NewExpression(w1, w2), kiwi.OpEq)
		err = solver2.AddConstraint(cnC)
		if err == nil || err.Error() != "unsatisfiable constraint" {
			t.Errorf("expected 'unsatisfiable constraint', got %v", err)
		}

		// GetConstraints check
		solver3 := kiwi.NewSolver()
		widthA := kiwi.NewVariable()
		widthB := kiwi.NewVariable()
		cn_1 := kiwi.NewConstraint(kiwi.NewExpression(widthA, 100), kiwi.OpEq)
		solver3.AddConstraint(cn_1)
		cn_2 := kiwi.NewConstraint(kiwi.NewExpression(widthB, 100), kiwi.OpEq)
		solver3.AddConstraint(cn_2)
		cns := solver3.GetConstraints()
		found1, found2 := false, false
		for _, c := range cns {
			if c == cn_1 {
				found1 = true
			}
			if c == cn_2 {
				found2 = true
			}
		}
		if !found1 || !found2 {
			t.Errorf("expected both constraints in GetConstraints()")
		}

		solver3.RemoveConstraint(cn_1)
		cns = solver3.GetConstraints()
		found1 = false
		for _, c := range cns {
			if c == cn_1 {
				found1 = true
			}
		}
		if found1 {
			t.Errorf("cn_1 should have been removed")
		}
	})

	t.Run("Constraint raw syntax suite", func(t *testing.T) {
		solver := kiwi.NewSolver()
		leftVar := kiwi.NewVariable()
		leftCn := kiwi.NewConstraint(kiwi.NewExpression([]any{-1, leftVar}, 10), kiwi.OpEq)
		solver.AddConstraint(leftCn)
		solver.UpdateVariables()
		if leftVar.Value() != 10 {
			t.Errorf("expected leftVar == 10, got %f", leftVar.Value())
		}

		widthVar := kiwi.NewVariable()
		solver.AddEditVariable(widthVar, kiwi.StrengthStrong)
		solver.SuggestValue(widthVar, 200)
		solver.UpdateVariables()
		if widthVar.Value() != 200 {
			t.Errorf("expected widthVar == 200, got %f", widthVar.Value())
		}

		rightVar := kiwi.NewVariable()
		rightCn := kiwi.NewConstraint(kiwi.NewExpression([]any{-1, rightVar}, leftVar, widthVar), kiwi.OpEq)
		solver.AddConstraint(rightCn)
		solver.UpdateVariables()
		if rightVar.Value() != 210 {
			t.Errorf("expected rightVar == 210, got %f", rightVar.Value())
		}

		centerXVar := kiwi.NewVariable()
		centerCn := kiwi.NewConstraint(kiwi.NewExpression([]any{-1, centerXVar}, leftVar, []any{0.5, widthVar}), kiwi.OpEq)
		solver.AddConstraint(centerCn)
		solver.UpdateVariables()
		if centerXVar.Value() != 110 {
			t.Errorf("expected centerXVar == 110, got %f", centerXVar.Value())
		}
	})

	t.Run("Constraint new syntax suite", func(t *testing.T) {
		solver := kiwi.NewSolver()
		left := kiwi.NewVariable()
		width := kiwi.NewVariable()
		top := kiwi.NewVariable()
		height := kiwi.NewVariable()
		right := kiwi.NewVariable()
		bottom := kiwi.NewVariable()
		centerX := kiwi.NewVariable()
		leftOfCenterX := kiwi.NewVariable()

		solver.AddEditVariable(left, kiwi.StrengthStrong)
		solver.AddEditVariable(width, kiwi.StrengthStrong)
		solver.AddEditVariable(top, kiwi.StrengthStrong)
		solver.AddEditVariable(height, kiwi.StrengthStrong)

		solver.SuggestValue(left, 0)
		solver.SuggestValue(width, 500)
		solver.SuggestValue(top, 0)
		solver.SuggestValue(height, 300)

		// right == left.plus(width) => 500
		solver.AddConstraint(kiwi.NewConstraint(right, kiwi.OpEq, left.Plus(width)))
		solver.UpdateVariables()
		if right.Value() != 500 {
			t.Errorf("expected right == 500, got %f", right.Value())
		}

		// centerX == left.plus(width.divide(2)) => 250
		solver.AddConstraint(kiwi.NewConstraint(centerX, kiwi.OpEq, left.Plus(width.Divide(2))))
		solver.UpdateVariables()
		if centerX.Value() != 250 {
			t.Errorf("expected centerX == 250, got %f", centerX.Value())
		}

		// leftOfCenterX == left.plus(width.divide(2)).minus(10) => 240
		solver.AddConstraint(kiwi.NewConstraint(leftOfCenterX, kiwi.OpEq, left.Plus(width.Divide(2)).Minus(10)))
		solver.UpdateVariables()
		if leftOfCenterX.Value() != 240 {
			t.Errorf("expected leftOfCenterX == 240, got %f", leftOfCenterX.Value())
		}

		// createConstraint(bottom, OpEq, top.plus(height)) => 300
		_, err := solver.CreateConstraint(bottom, kiwi.OpEq, top.Plus(height))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		solver.UpdateVariables()
		if bottom.Value() != 300 {
			t.Errorf("expected bottom == 300, got %f", bottom.Value())
		}
	})
}
