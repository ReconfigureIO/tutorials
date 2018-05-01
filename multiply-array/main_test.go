package main

import (
	"testing"
	"testing/quick"
)

func TestMultiply(t *testing.T) {
	// Check that our Multiply function multiplies its input by 2
	f := func(x uint32) bool {
		result := Multiply(x)
		return (result == (x * 2))
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
