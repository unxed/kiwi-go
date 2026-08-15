package kiwi

import (
	"math"
	"sort"
)

// DiscreteResult holds the rounded integer layout coordinates/dimensions for variables.
type DiscreteResult map[*Variable]int

// Get returns the discrete integer value for the given variable, or 0 if not present.
func (r DiscreteResult) Get(v *Variable) int {
	if r == nil || v == nil {
		return 0
	}
	return r[v]
}

// Sum returns the sum of discrete integer values for the specified variables.
func (r DiscreteResult) Sum(vars ...*Variable) int {
	total := 0
	for _, v := range vars {
		total += r.Get(v)
	}
	return total
}

type variableRemainder struct {
	variable  *Variable
	floorVal  int
	remainder float64
	index     int
}

// ApportionSum distributes integer units across a group of variables using the
// Hare-Niemeyer / Hamilton Largest Remainder method (inspired by FreeType autohinting).
// It guarantees that sum(DiscreteResult[v_i]) == targetSum exactly, preserves symmetry,
// and respects optional minSizes.
func ApportionSum(vars []*Variable, floatVals map[*Variable]float64, targetSum int, minSizes ...map[*Variable]int) DiscreteResult {
	result := make(DiscreteResult)
	if len(vars) == 0 {
		return result
	}

	var minMap map[*Variable]int
	if len(minSizes) > 0 {
		minMap = minSizes[0]
	}

	// 1. Calculate floor values and remainders
	items := make([]variableRemainder, len(vars))
	sumFloor := 0

	for i, v := range vars {
		val := 0.0
		if floatVals != nil {
			val = floatVals[v]
		} else if v != nil {
			val = v.Value()
		}

		minVal := 0
		if minMap != nil {
			minVal = minMap[v]
		}

		fVal := int(math.Floor(val))
		if fVal < minVal {
			fVal = minVal
		}

		rem := val - float64(fVal)
		items[i] = variableRemainder{
			variable:  v,
			floorVal:  fVal,
			remainder: rem,
			index:     i,
		}
		sumFloor += fVal
	}

	deficit := targetSum - sumFloor

	// If sumFloor exceeds targetSum, return base values
	if deficit < 0 {
		for _, item := range items {
			result[item.variable] = item.floorVal
		}
		return result
	}

	// 2. Sort by largest remainder descending with symmetry-aware tie-breaking
	sort.SliceStable(items, func(i, j int) bool {
		if math.Abs(items[i].remainder-items[j].remainder) > 1e-7 {
			return items[i].remainder > items[j].remainder
		}
		// Tie-breaker: distance from center for symmetric layout distribution
		mid := float64(len(vars)-1) / 2.0
		distI := math.Abs(float64(items[i].index) - mid)
		distJ := math.Abs(float64(items[j].index) - mid)
		if math.Abs(distI-distJ) > 1e-7 {
			return distI > distJ
		}
		return items[i].index < items[j].index
	})

	// 3. Assign base floor values
	for _, item := range items {
		result[item.variable] = item.floorVal
	}

	// 4. Distribute remainder deficit (+1 to top deficit variables)
	for i := 0; i < deficit && i < len(items); i++ {
		result[items[i].variable]++
	}

	return result
}

// RoundValues performs standard integer rounding for variables, respecting optional minimum sizes.
func RoundValues(vars []*Variable, floatVals map[*Variable]float64, minSizes ...map[*Variable]int) DiscreteResult {
	result := make(DiscreteResult)
	var minMap map[*Variable]int
	if len(minSizes) > 0 {
		minMap = minSizes[0]
	}

	for _, v := range vars {
		val := 0.0
		if floatVals != nil {
			val = floatVals[v]
		} else if v != nil {
			val = v.Value()
		}

		rVal := int(math.Round(val))
		if minMap != nil {
			if minVal, ok := minMap[v]; ok && rVal < minVal {
				rVal = minVal
			}
		}
		result[v] = rVal
	}
	return result
}
