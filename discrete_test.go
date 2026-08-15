package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestApportionSumEqualColumns(t *testing.T) {
	v1 := kiwi.NewVariable("col1")
	v2 := kiwi.NewVariable("col2")
	v3 := kiwi.NewVariable("col3")

	vars := []*kiwi.Variable{v1, v2, v3}
	floatVals := map[*kiwi.Variable]float64{
		v1: 26.66666667,
		v2: 26.66666667,
		v3: 26.66666667,
	}

	// Target 80 chars
	res := kiwi.ApportionSum(vars, floatVals, 80)
	if res.Sum(vars...) != 80 {
		t.Errorf("expected sum 80, got %d", res.Sum(vars...))
	}

	if res.Get(v1) < 26 || res.Get(v1) > 27 {
		t.Errorf("unexpected value for v1: %d", res.Get(v1))
	}
	if res.Get(v2) < 26 || res.Get(v2) > 27 {
		t.Errorf("unexpected value for v2: %d", res.Get(v2))
	}
	if res.Get(v3) < 26 || res.Get(v3) > 27 {
		t.Errorf("unexpected value for v3: %d", res.Get(v3))
	}
}

func TestApportionSumMinSize(t *testing.T) {
	v1 := kiwi.NewVariable("a")
	v2 := kiwi.NewVariable("b")

	vars := []*kiwi.Variable{v1, v2}
	floatVals := map[*kiwi.Variable]float64{
		v1: 0.2,
		v2: 9.8,
	}

	minSizes := map[*kiwi.Variable]int{
		v1: 1, // Minimum 1 char width
	}

	res := kiwi.ApportionSum(vars, floatVals, 10, minSizes)
	if res.Get(v1) < 1 {
		t.Errorf("expected v1 >= 1, got %d", res.Get(v1))
	}
	if res.Sum(vars...) != 10 {
		t.Errorf("expected total sum 10, got %d", res.Sum(vars...))
	}
}

func TestRoundValuesAndResultHelpers(t *testing.T) {
	v1 := kiwi.NewVariable("x")
	v2 := kiwi.NewVariable("y")

	v1.SetValue(10.4)
	v2.SetValue(20.6)

	res := kiwi.RoundValues([]*kiwi.Variable{v1, v2}, nil)
	if res.Get(v1) != 10 {
		t.Errorf("expected 10, got %d", res.Get(v1))
	}
	if res.Get(v2) != 21 {
		t.Errorf("expected 21, got %d", res.Get(v2))
	}
	if res.Sum(v1, v2) != 31 {
		t.Errorf("expected sum 31, got %d", res.Sum(v1, v2))
	}

	// Nil safety checks
	var nilRes kiwi.DiscreteResult
	if nilRes.Get(v1) != 0 || nilRes.Sum(v1) != 0 {
		t.Errorf("nil DiscreteResult should return 0")
	}

	// Empty slice apportion check
	emptyRes := kiwi.ApportionSum(nil, nil, 10)
	if len(emptyRes) != 0 {
		t.Errorf("expected empty result for nil vars")
	}
}
