package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestApportionSumNegativeDeficitExcess(t *testing.T) {
	v1 := kiwi.NewVariable("v1")
	v2 := kiwi.NewVariable("v2")

	// Continuous values sum to 10.2, but target is 9
	vars := []*kiwi.Variable{v1, v2}
	floatVals := map[*kiwi.Variable]float64{
		v1: 5.1,
		v2: 5.1,
	}

	res := kiwi.ApportionSum(vars, floatVals, 9)
	if res.Sum(vars...) != 9 {
		t.Errorf("expected sum 9, got %d", res.Sum(vars...))
	}
}

func TestApportionSumNegativeCoordinates(t *testing.T) {
	v1 := kiwi.NewVariable("offscreen_x")
	v2 := kiwi.NewVariable("margin_x")

	vars := []*kiwi.Variable{v1, v2}
	floatVals := map[*kiwi.Variable]float64{
		v1: -10.3,
		v2: -5.7,
	}

	res := kiwi.ApportionSum(vars, floatVals, -16)
	if res.Sum(vars...) != -16 {
		t.Errorf("expected sum -16, got %d", res.Sum(vars...))
	}
}

func TestApportionSumNilVariableFilter(t *testing.T) {
	v1 := kiwi.NewVariable("v1")
	vars := []*kiwi.Variable{v1, nil}

	res := kiwi.ApportionSum(vars, nil, 10)
	if res.Get(v1) != 10 {
		t.Errorf("expected v1 == 10, got %d", res.Get(v1))
	}
	if res.Get(nil) != 0 {
		t.Errorf("nil variable should evaluate to 0")
	}
}
func TestApportionSumMultiUnitExcess(t *testing.T) {
	v1 := kiwi.NewVariable("v1")
	v2 := kiwi.NewVariable("v2")

	vars := []*kiwi.Variable{v1, v2}
	// Floor values 10 and 10 (sum 20), target is 15 (excess = 5 > len(vars))
	floatVals := map[*kiwi.Variable]float64{
		v1: 10.2,
		v2: 10.2,
	}

	res := kiwi.ApportionSum(vars, floatVals, 15)
	if res.Sum(vars...) != 15 {
		t.Errorf("expected sum 15, got %d", res.Sum(vars...))
	}
}
