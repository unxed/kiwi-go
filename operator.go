package kiwi

import "fmt"

// Operator defines the linear constraint relational operators.
type Operator int

const (
	// OpLe represents Less than or equal (<=).
	OpLe Operator = 0
	// OpGe represents Greater than or equal (>=).
	OpGe Operator = 1
	// OpEq represents Equal (==).
	OpEq Operator = 2
)

// Short aliases matching TypeScript enum names
const (
	Le = OpLe
	Ge = OpGe
	Eq = OpEq
)

func (op Operator) String() string {
	switch op {
	case OpLe:
		return "<="
	case OpGe:
		return ">="
	case OpEq:
		return "="
	default:
		return fmt.Sprintf("Operator(%d)", int(op))
	}
}
