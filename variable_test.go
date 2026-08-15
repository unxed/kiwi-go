package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestVariable(t *testing.T) {
	v1 := kiwi.NewVariable()
	if v1.Value() != 0 {
		t.Errorf("expected initial value 0, got %f", v1.Value())
	}

	v2 := kiwi.NewVariable("somename")
	if v2.Name() != "somename" {
		t.Errorf("expected name 'somename', got %s", v2.Name())
	}

	v2.SetName("skiwi")
	if v2.Name() != "skiwi" {
		t.Errorf("expected name 'skiwi', got %s", v2.Name())
	}

	if v1.ID() == v2.ID() {
		t.Errorf("expected unique IDs, got %d and %d", v1.ID(), v2.ID())
	}

	var calledVal, calledPrev float64
	callCount := 0
	v1.Subscribe(func(val, prev float64) {
		callCount++
		calledVal = val
		calledPrev = prev
	})

	v1.SetValue(200)
	if callCount != 1 || calledVal != 200 || calledPrev != 0 {
		t.Errorf("expected callback with 200, 0, got count %d, val %f, prev %f", callCount, calledVal, calledPrev)
	}

	// Setting same value should not trigger callback
	v1.SetValue(200)
	if callCount != 1 {
		t.Errorf("expected callCount 1, got %d", callCount)
	}

	v1.Unsubscribe()
	v1.SetValue(300)
	if callCount != 1 {
		t.Errorf("expected callCount 1 after unsubscribe, got %d", callCount)
	}

	v1.SetContext("testCtx")
	if v1.Context() != "testCtx" {
		t.Errorf("expected context 'testCtx', got %v", v1.Context())
	}

	str := v1.String()
	if str == "" {
		t.Errorf("expected non-empty string representation")
	}
}
func TestVariableStringRepresentation(t *testing.T) {
	v := kiwi.NewVariable("x")
	v.SetContext("myContext")
	v.SetValue(42)

	str := v.String()
	expected := "myContext[x:42]"
	if str != expected {
		t.Errorf("expected %s, got %s", expected, str)
	}
}
