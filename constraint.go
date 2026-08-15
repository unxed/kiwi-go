package kiwi

import (
	"fmt"
	"sync/atomic"
)

var globalCnID int64

// Constraint represents a linear constraint equation.
type Constraint struct {
	id         int
	expression *Expression
	operator   Operator
	strength   float64
}

// NewConstraint creates a new Constraint instance.
// Arguments: (expression|variable, operator, [rhs], [strength])
func NewConstraint(expression any, operator Operator, rhsAndStrength ...any) *Constraint {
	var rhs any = nil
	strength := StrengthRequired

	if len(rhsAndStrength) > 0 {
		rhs = rhsAndStrength[0]
	}
	if len(rhsAndStrength) > 1 {
		if st, ok := toFloat(rhsAndStrength[1]); ok {
			strength = st
		}
	}

	st := ClipStrength(strength)
	id := atomic.AddInt64(&globalCnID, 1) - 1

	var expr *Expression
	if rhs == nil {
		if e, ok := expression.(*Expression); ok {
			expr = e
		} else if v, ok := expression.(*Variable); ok {
			expr = NewExpression(v)
		} else {
			expr = NewExpression(expression)
		}
	} else {
		if e, ok := expression.(*Expression); ok {
			expr = e.Minus(rhs)
		} else if v, ok := expression.(*Variable); ok {
			expr = v.Minus(rhs)
		} else {
			expr = NewExpression(expression).Minus(rhs)
		}
	}

	return &Constraint{
		id:         int(id),
		expression: expr,
		operator:   operator,
		strength:   st,
	}
}

// ID returns the unique id number of the constraint.
func (c *Constraint) ID() int {
	return c.id
}

// Expression returns the expression of the constraint.
func (c *Constraint) Expression() *Expression {
	return c.expression
}

// Op returns the relational operator of the constraint.
func (c *Constraint) Op() Operator {
	return c.operator
}

// Strength returns the strength of the constraint.
func (c *Constraint) Strength() float64 {
	return c.strength
}

func (c *Constraint) String() string {
	opStr := ["<=", ">=", "="][c.operator]
	return fmt.Sprintf("%s %s 0 (%g)", c.expression.String(), opStr, c.strength)
}
