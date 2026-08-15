package kiwi

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var globalVarID int64

// VariableCallback is a subscription function called whenever the variable's value changes.
type VariableCallback func(value, previousValue float64)

// Variable is the primary user constraint variable.
type Variable struct {
	mu       sync.RWMutex
	id       int
	name     string
	val      float64
	context  any
	callback VariableCallback
}

// NewVariable constructs a new Variable with an optional name.
func NewVariable(name ...string) *Variable {
	n := ""
	if len(name) > 0 {
		n = name[0]
	}
	id := atomic.AddInt64(&globalVarID, 1) - 1
	return &Variable{
		id:   int(id),
		name: n,
	}
}

// ID returns the unique id number of the variable.
func (v *Variable) ID() int {
	return v.id
}

// Name returns the name of the variable.
func (v *Variable) Name() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.name
}

// SetName sets the name of the variable.
func (v *Variable) SetName(name string) {
	v.mu.Lock()
	v.name = name
	v.mu.Unlock()
}

// Context returns the user context object of the variable.
func (v *Variable) Context() any {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.context
}

// SetContext sets the user context object of the variable.
func (v *Variable) SetContext(context any) {
	v.mu.Lock()
	v.context = context
	v.mu.Unlock()
}

// Value returns the calculated value of the variable.
func (v *Variable) Value() float64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.val
}

// SetValue sets the value of the variable and notifies subscribers if changed.
func (v *Variable) SetValue(value float64) {
	v.mu.Lock()
	prev := v.val
	v.val = value
	cb := v.callback
	v.mu.Unlock()

	if cb != nil && prev != value {
		cb(value, prev)
	}
}

// Subscribe sets a callback for whenever the value changes.
func (v *Variable) Subscribe(callback VariableCallback) {
	v.mu.Lock()
	v.callback = callback
	v.mu.Unlock()
}

// Unsubscribe stops the variable from calling the callback when the value changes.
func (v *Variable) Unsubscribe() {
	v.mu.Lock()
	v.callback = nil
	v.mu.Unlock()
}

// Plus creates a new Expression by adding a number, variable or expression.
func (v *Variable) Plus(value any) *Expression {
	return NewExpression(v, value)
}

// Minus creates a new Expression by subtracting a number, variable or expression.
func (v *Variable) Minus(value any) *Expression {
	if val, ok := toFloat(value); ok {
		return NewExpression(v, -val)
	}
	return NewExpression(v, []any{-1.0, value})
}

// Multiply creates a new Expression by multiplying with a fixed number.
func (v *Variable) Multiply(coefficient float64) *Expression {
	return NewExpression([]any{coefficient, v})
}

// Divide creates a new Expression by dividing with a fixed number.
func (v *Variable) Divide(coefficient float64) *Expression {
	return NewExpression([]any{1.0 / coefficient, v})
}

// String returns string representation of the variable.
func (v *Variable) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return fmt.Sprintf("%v[%s:%g]", v.context, v.name, v.val)
}

// MarshalJSON returns the JSON representation of the variable.
func (v *Variable) MarshalJSON() ([]byte, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return []byte(fmt.Sprintf(`{"name":%q,"value":%g}`, v.name, v.val)), nil
}
