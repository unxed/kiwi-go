package kiwi_test

import (
	"fmt"

	"github.com/unxed/kiwi-go"
)

func ExampleSolver() {
	solver := kiwi.NewSolver()

	left := kiwi.NewVariable("left")
	width := kiwi.NewVariable("width")
	right := kiwi.NewVariable("right")

	_ = solver.AddEditVariable(left, kiwi.StrengthStrong)
	_ = solver.AddEditVariable(width, kiwi.StrengthStrong)

	_ = solver.SuggestValue(left, 100)
	_ = solver.SuggestValue(width, 400)

	// right == left + width
	cn := kiwi.NewConstraint(right, kiwi.OpEq, left.Plus(width))
	_ = solver.AddConstraint(cn)

	solver.UpdateVariables()

	fmt.Printf("right = %.0f\n", right.Value())

	// Output:
	// right = 500
}

func ExampleVariable() {
	x := kiwi.NewVariable("x")
	expr := x.Multiply(2).Plus(10)

	x.SetValue(15)
	fmt.Printf("expression value = %.0f\n", expr.Value())

	// Output:
	// expression value = 40
}
