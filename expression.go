package kiwi

import (
	"errors"
	"fmt"
	"strings"
)

// Expression represents an expression of variable terms and a constant.
type Expression struct {
	terms    *IndexedMap[*Variable, float64]
	constant float64
}

// NewExpression creates a new Expression summing the given arguments.
// Arguments can be numbers (int/float), *Variable, *Expression, or 2-item slices []any{coeff, varOrExpr}.
func NewExpression(args ...any) *Expression {
	terms, constant := parseArgs(args)
	return &Expression{
		terms:    terms,
		constant: constant,
	}
}

// Terms returns the mapping of terms in the expression.
func (e *Expression) Terms() *IndexedMap[*Variable, float64] {
	return e.terms
}

// Constant returns the constant of the expression.
func (e *Expression) Constant() float64 {
	return e.constant
}

// Value returns the computed value of the expression.
func (e *Expression) Value() float64 {
	result := e.constant
	for i := 0; i < e.terms.Size(); i++ {
		pair := e.terms.ItemAt(i)
		result += pair.First.Value() * pair.Second
	}
	return result
}

// Plus creates a new Expression by adding a number, variable, or expression.
func (e *Expression) Plus(value any) *Expression {
	return NewExpression(e, value)
}

// Minus creates a new Expression by subtracting a number, variable, or expression.
func (e *Expression) Minus(value any) *Expression {
	if val, ok := toFloat(value); ok {
		return NewExpression(e, -val)
	}
	return NewExpression(e, []any{-1.0, value})
}

// Multiply creates a new Expression by multiplying with a fixed number.
func (e *Expression) Multiply(coefficient float64) *Expression {
	return NewExpression([]any{coefficient, e})
}

// Divide creates a new Expression by dividing with a fixed number.
func (e *Expression) Divide(coefficient float64) *Expression {
	return NewExpression([]any{1.0 / coefficient, e})
}

// IsConstant returns true if the expression contains no variable terms.
func (e *Expression) IsConstant() bool {
	return e.terms.Size() == 0
}

func (e *Expression) String() string {
	var parts []string
	arr := e.terms.Array()
	for _, pair := range arr {
		parts = append(parts, fmt.Sprintf("%g*%s", pair.Second, pair.First.String()))
	}
	result := strings.Join(parts, " + ")

	if !e.IsConstant() && e.constant != 0 {
		result += " + "
	}
	if e.IsConstant() || e.constant != 0 {
		result += fmt.Sprintf("%g", e.constant)
	}
	return result
}

func parseArgs(args []any) (*IndexedMap[*Variable, float64], float64) {
	constant := 0.0
	factory := func() float64 { return 0.0 }
	terms := NewIndexedMap[*Variable, float64]()

	for _, item := range args {
		if val, ok := toFloat(item); ok {
			constant += val
			continue
		}

		switch v := item.(type) {
		case *Variable:
			terms.SetDefault(v, factory).Second += 1.0
		case *Expression:
			constant += v.Constant()
			t2 := v.Terms()
			for j := 0; j < t2.Size(); j++ {
				termPair := t2.ItemAt(j)
				terms.SetDefault(termPair.First, factory).Second += termPair.Second
			}
		case []any:
			if len(v) != 2 {
				panic(errors.New("array must have length 2"))
			}
			coeff, ok := toFloat(v[0])
			if !ok {
				panic(errors.New("array item 0 must be a number"))
			}
			target := v[1]
			switch t := target.(type) {
			case *Variable:
				terms.SetDefault(t, factory).Second += coeff
			case *Expression:
				constant += t.Constant() * coeff
				t2 := t.Terms()
				for j := 0; j < t2.Size(); j++ {
					termPair := t2.ItemAt(j)
					terms.SetDefault(termPair.First, factory).Second += termPair.Second * coeff
				}
			default:
				panic(errors.New("array item 1 must be a variable or expression"))
			}
		default:
			panic(fmt.Errorf("invalid Expression argument: %v", item))
		}
	}
	return terms, constant
}

func toFloat(item any) (float64, bool) {
	switch v := item.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	default:
		return 0, false
	}
}
