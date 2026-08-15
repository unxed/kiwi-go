package kiwi_test

import (
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestStrength(t *testing.T) {
	required := kiwi.CreateStrength(1000, 1000, 1000)
	if required != kiwi.StrengthRequired {
		t.Errorf("expected StrengthRequired %f, got %f", kiwi.StrengthRequired, required)
	}

	strong := kiwi.CreateStrength(1, 0, 0)
	if strong != kiwi.StrengthStrong {
		t.Errorf("expected StrengthStrong %f, got %f", kiwi.StrengthStrong, strong)
	}

	medium := kiwi.CreateStrength(0, 1, 0)
	if medium != kiwi.StrengthMedium {
		t.Errorf("expected StrengthMedium %f, got %f", kiwi.StrengthMedium, medium)
	}

	weak := kiwi.CreateStrength(0, 0, 1)
	if weak != kiwi.StrengthWeak {
		t.Errorf("expected StrengthWeak %f, got %f", kiwi.StrengthWeak, weak)
	}

	clipped := kiwi.ClipStrength(2000000000)
	if clipped != kiwi.StrengthRequired {
		t.Errorf("expected clipped value to be StrengthRequired, got %f", clipped)
	}

	clippedNeg := kiwi.ClipStrength(-100)
	if clippedNeg != 0 {
		t.Errorf("expected clipped negative value to be 0, got %f", clippedNeg)
	}

	var s kiwi.Strength
	if s.Required() != kiwi.StrengthRequired {
		t.Errorf("expected Strength.Required() to match StrengthRequired")
	}
	if s.Strong() != kiwi.StrengthStrong {
		t.Errorf("expected Strength.Strong() to match StrengthStrong")
	}
	if s.Medium() != kiwi.StrengthMedium {
		t.Errorf("expected Strength.Medium() to match StrengthMedium")
	}
	if s.Weak() != kiwi.StrengthWeak {
		t.Errorf("expected Strength.Weak() to match StrengthWeak")
	}
	if s.Clip(500) != 500 {
		t.Errorf("expected Strength.Clip(500) to return 500")
	}
}
