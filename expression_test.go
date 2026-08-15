package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestExpressionConstructors(t *testing.T) {
	e1 := kiwi.NewExpression()
	if e1 == nil {
		t.Fatal("expected non-nil expression")
	}

	v1 := kiwi.NewVariable("x")
	e2 := kiwi.NewExpression(v1)
	if e2.Terms().Size() != 1 {
		t.Errorf("expected 1 term, got %d", e2.Terms().Size())
	}

	e3 := kiwi.NewExpression([]any{-1, v1})
	p, ok := e3.Terms().Find(v1)
	if !ok || p.Second != -1 {
		t.Errorf("expected coefficient -1 for v1")
	}
}

func TestExpressionErrors(t *testing.T) {
	assertPanic := func(fn func(), expected string) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected panic '%s', got nil", expected)
				return
			}
			var msg string
			if err, ok := r.(error); ok {
				msg = err.Error()
			} else {
				msg = r.(string)
			}
			if msg != expected {
				t.Errorf("expected panic '%s', got '%s'", expected, msg)
			}
		}()
		fn()
	}

	v := kiwi.NewVariable()
	assertPanic(func() {
		kiwi.NewExpression([]any{v, -1})
	}, "array item 0 must be a number")

	assertPanic(func() {
		kiwi.NewExpression([]any{-1, 100})
	}, "array item 1 must be a variable or expression")

	assertPanic(func() {
		kiwi.NewExpression([]any{-1})
	}, "array must have length 2")

	assertPanic(func() {
		kiwi.NewExpression("invalid_arg")
	}, "invalid Expression argument: invalid_arg")
}

func TestExpressionConstants(t *testing.T) {
	e1 := kiwi.NewExpression(10, 20, 30, 40)
	if e1.Constant() != 100 {
		t.Errorf("expected constant 100, got %f", e1.Constant())
	}

	e2 := kiwi.NewExpression([]any{-1, kiwi.NewExpression(10)})
	if e2.Constant() != -10 {
		t.Errorf("expected constant -10, got %f", e2.Constant())
	}

	e3 := kiwi.NewExpression(kiwi.NewExpression(10), kiwi.NewExpression(20))
	if e3.Constant() != 30 {
		t.Errorf("expected constant 30, got %f", e3.Constant())
	}

	e4 := kiwi.NewExpression(20, []any{0.5, kiwi.NewExpression(10)}, -10)
	if e4.Constant() != 15 {
		t.Errorf("expected constant 15, got %f", e4.Constant())
	}
}

func TestExpressionArithmetic(t *testing.T) {
	v := kiwi.NewVariable("x")
	v.SetValue(10)

	expr := v.Plus(5)
	if expr.Value() != 15 {
		t.Errorf("expected 15, got %f", expr.Value())
	}

	expr2 := expr.Minus(3)
	if expr2.Value() != 12 {
		t.Errorf("expected 12, got %f", expr2.Value())
	}

	expr3 := v.Multiply(3)
	if expr3.Value() != 30 {
		t.Errorf("expected 30, got %f", expr3.Value())
	}

	expr4 := v.Divide(2)
	if expr4.Value() != 5 {
		t.Errorf("expected 5, got %f", expr4.Value())
	}

	if expr.IsConstant() {
		t.Errorf("expected expr with variable term to not be constant")
	}

	cExpr := kiwi.NewExpression(42)
	if !cExpr.IsConstant() {
		t.Errorf("expected cExpr to be constant")
	}

	str := expr.String()
	if str == "" {
		t.Errorf("expected non-empty string representation")
	}
}
