package main

import (
	"os"
	"testing"
)

func TestDemoExecution(t *testing.T) {
	os.Setenv("KIWI_TESTING", "1")
	defer os.Unsetenv("KIWI_TESTING")
	main()
}
