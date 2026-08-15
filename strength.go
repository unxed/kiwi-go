package kiwi

import "math"

// CreateStrength creates a new symbolic strength value from strong (a), medium (b),
// weak (c) components, and an optional weight w (defaulting to 1.0).
func CreateStrength(a, b, c float64, weight ...float64) float64 {
	w := 1.0
	if len(weight) > 0 {
		w = weight[0]
	}
	var result float64
	result += math.Max(0.0, math.Min(1000.0, a*w)) * 1000000.0
	result += math.Max(0.0, math.Min(1000.0, b*w)) * 1000.0
	result += math.Max(0.0, math.Min(1000.0, c*w))
	return result
}

var (
	// StrengthRequired is the 'required' symbolic strength.
	StrengthRequired = CreateStrength(1000.0, 1000.0, 1000.0)
	// StrengthStrong is the 'strong' symbolic strength.
	StrengthStrong = CreateStrength(1.0, 0.0, 0.0)
	// StrengthMedium is the 'medium' symbolic strength.
	StrengthMedium = CreateStrength(0.0, 1.0, 0.0)
	// StrengthWeak is the 'weak' symbolic strength.
	StrengthWeak = CreateStrength(0.0, 0.0, 1.0)
)

// ClipStrength clips a symbolic strength to the allowed min and max.
func ClipStrength(value float64) float64 {
	return math.Max(0.0, math.Min(StrengthRequired, value))
}

// Strength provides a namespace matching the TypeScript Strength class API.
type Strength struct{}

func (Strength) Create(a, b, c float64, weight ...float64) float64 {
	return CreateStrength(a, b, c, weight...)
}

func (Strength) Required() float64 {
	return StrengthRequired
}

func (Strength) Strong() float64 {
	return StrengthStrong
}

func (Strength) Medium() float64 {
	return StrengthMedium
}

func (Strength) Weak() float64 {
	return StrengthWeak
}

func (Strength) Clip(value float64) float64 {
	return ClipStrength(value)
}
