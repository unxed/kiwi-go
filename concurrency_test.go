package kiwi_test

import (
	"sync"
	"testing"

	"github.com/unxed/kiwi-go"
)

func TestVariableConcurrency(t *testing.T) {
	v := kiwi.NewVariable("x")
	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(val float64) {
			defer wg.Done()
			v.SetValue(val)
		}(float64(i))
	}

	// Concurrent readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = v.Value()
			_ = v.Name()
			_ = v.String()
			_, _ = v.MarshalJSON()
		}()
	}

	wg.Wait()
}
