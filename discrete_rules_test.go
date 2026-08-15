package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestSnapToGrid(t *testing.T) {
	v1 := kiwi.NewVariable("cjk1")
	v2 := kiwi.NewVariable("cjk2")

	res := kiwi.DiscreteResult{
		v1: 5,
		v2: 3,
	}

	hinter := kiwi.NewRuleHinter()
	hinter.AddDirective(
		kiwi.SnapToGrid(v1, 2),
		kiwi.SnapToGrid(v2, 2),
	)
	hinter.Apply(res)

	if res.Get(v1) != 6 {
		t.Errorf("expected 5 snapped to 6, got %d", res.Get(v1))
	}
	if res.Get(v2) != 2 {
		t.Errorf("expected 3 snapped to 2, got %d", res.Get(v2))
	}
}

func TestClampMinMax(t *testing.T) {
	v1 := kiwi.NewVariable("min_test")
	v2 := kiwi.NewVariable("max_test")
	v3 := kiwi.NewVariable("ok_test")

	res := kiwi.DiscreteResult{
		v1: 2,
		v2: 50,
		v3: 15,
	}

	hinter := kiwi.NewRuleHinter()
	hinter.AddDirective(
		kiwi.ClampMinMax(v1, 10, 20),
		kiwi.ClampMinMax(v2, 10, 20),
		kiwi.ClampMinMax(v3, 10, 20),
	)
	hinter.Apply(res)

	if res.Get(v1) != 10 {
		t.Errorf("expected v1 clamped to min 10, got %d", res.Get(v1))
	}
	if res.Get(v2) != 20 {
		t.Errorf("expected v2 clamped to max 20, got %f", float64(res.Get(v2)))
	}
	if res.Get(v3) != 15 {
		t.Errorf("expected v3 unchanged 15, got %d", res.Get(v3))
	}
}

func TestEqualizeGroup(t *testing.T) {
	v1 := kiwi.NewVariable("col1")
	v2 := kiwi.NewVariable("col2")
	v3 := kiwi.NewVariable("col3")

	res := kiwi.DiscreteResult{
		v1: 10,
		v2: 11,
		v3: 12,
	}

	hinter := kiwi.NewRuleHinter()
	hinter.AddDirective(kiwi.EqualizeGroup(v1, v2, v3))
	hinter.Apply(res)

	if res.Get(v1) != 11 || res.Get(v2) != 11 || res.Get(v3) != 11 {
		t.Errorf("expected all equalized to 11, got %d, %d, %d", res.Get(v1), res.Get(v2), res.Get(v3))
	}
}

func TestAlignEdgesAndCustomDirective(t *testing.T) {
	v1 := kiwi.NewVariable("src")
	v2 := kiwi.NewVariable("dst")
	v3 := kiwi.NewVariable("custom")

	res := kiwi.DiscreteResult{
		v1: 100,
		v2: 10,
		v3: 5,
	}

	customCalled := false
	hinter := kiwi.NewRuleHinter()
	hinter.AddDirective(
		kiwi.AlignEdges(v2, v1),
		kiwi.CustomDirective(func(r kiwi.DiscreteResult) {
			customCalled = true
			r[v3] = 42
		}),
	)
	hinter.Apply(res)

	if res.Get(v2) != 100 {
		t.Errorf("expected v2 aligned to v1 (100), got %d", res.Get(v2))
	}
	if !customCalled || res.Get(v3) != 42 {
		t.Errorf("expected custom directive applied, got %d", res.Get(v3))
	}

	// Nil safety check
	hinter.Apply(nil)
}
