package kiwi

import "math"

// HintDirective defines an explicit hinting directive applied to discrete layout results
// (inspired by TrueType bytecode hinting instructions).
type HintDirective interface {
	Apply(res DiscreteResult)
}

// HintDirectiveFunc is a function adapter implementing HintDirective.
type HintDirectiveFunc func(res DiscreteResult)

// Apply invokes the underlying function.
func (f HintDirectiveFunc) Apply(res DiscreteResult) {
	if f != nil {
		f(res)
	}
}

// RuleHinter manages and executes a pipeline of explicit HintDirectives.
type RuleHinter struct {
	directives []HintDirective
}

// NewRuleHinter creates a new RuleHinter instance.
func NewRuleHinter() *RuleHinter {
	return &RuleHinter{
		directives: make([]HintDirective, 0),
	}
}

// AddDirective adds one or more HintDirectives to the hinter pipeline.
func (h *RuleHinter) AddDirective(directives ...HintDirective) *RuleHinter {
	h.directives = append(h.directives, directives...)
	return h
}

// Apply executes all directives sequentially on the DiscreteResult.
func (h *RuleHinter) Apply(res DiscreteResult) {
	if res == nil {
		return
	}
	for _, d := range h.directives {
		if d != nil {
			d.Apply(res)
		}
	}
}

// SnapToGrid snaps the variable's discrete value to the nearest multiple of step
// (e.g. step=2 for double-width character or grid-alignment constraints).
func SnapToGrid(v *Variable, step int) HintDirective {
	return HintDirectiveFunc(func(res DiscreteResult) {
		if v == nil || step <= 1 {
			return
		}
		val := res.Get(v)
		rem := val % step
		if rem != 0 {
			if rem*2 >= step {
				res[v] = val + (step - rem)
			} else {
				res[v] = val - rem
			}
		}
	})
}

// ClampMinMax clamps the variable's discrete value to the range [minVal, maxVal].
func ClampMinMax(v *Variable, minVal, maxVal int) HintDirective {
	return HintDirectiveFunc(func(res DiscreteResult) {
		if v == nil {
			return
		}
		val := res.Get(v)
		if val < minVal {
			res[v] = minVal
		} else if maxVal >= minVal && val > maxVal {
			res[v] = maxVal
		}
	})
}

// EqualizeGroup forces all variables in the group to have the exact same integer value (the rounded average).
func EqualizeGroup(vars ...*Variable) HintDirective {
	return HintDirectiveFunc(func(res DiscreteResult) {
		if len(vars) <= 1 {
			return
		}
		sum := 0
		count := 0
		for _, v := range vars {
			if v != nil {
				sum += res.Get(v)
				count++
			}
		}
		if count == 0 {
			return
		}
		avg := int(math.Round(float64(sum) / float64(count)))
		for _, v := range vars {
			if v != nil {
				res[v] = avg
			}
		}
	})
}

// AlignEdges sets variable v's discrete coordinate equal to target's discrete coordinate.
func AlignEdges(v *Variable, target *Variable) HintDirective {
	return HintDirectiveFunc(func(res DiscreteResult) {
		if v == nil || target == nil {
			return
		}
		res[v] = res.Get(target)
	})
}

// CustomDirective allows attaching a custom user-defined hinting callback function.
func CustomDirective(fn func(res DiscreteResult)) HintDirective {
	return HintDirectiveFunc(fn)
}
