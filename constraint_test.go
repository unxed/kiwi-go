package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestConstraintConstructors(t *testing.T) {
	v1 := kiwi.NewVariable("v1")
	v2 := kiwi.NewVariable("v2")

	cn1 := kiwi.NewConstraint(v1, kiwi.OpEq)
	if cn1.Op() != kiwi.OpEq {
		t.Errorf("expected OpEq, got %v", cn1.Op())
	}
	if cn1.Strength() != kiwi.StrengthRequired {
		t.Errorf("expected StrengthRequired, got %f", cn1.Strength())
	}

	cn2 := kiwi.NewConstraint(v1, kiwi.OpGe, v2, kiwi.StrengthMedium)
	if cn2.Op() != kiwi.OpGe {
		t.Errorf("expected OpGe, got %v", cn2.Op())
	}
	if cn2.Strength() != kiwi.StrengthMedium {
		t.Errorf("expected StrengthMedium, got %f", cn2.Strength())
	}

	if cn1.ID() == cn2.ID() {
		t.Errorf("expected unique constraint IDs")
	}

	str := cn2.String()
	if str == "" {
		t.Errorf("expected non-empty string representation")
	}
}
