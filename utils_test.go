package main

import (
	"testing"
)

func TestIntersectionOf(t *testing.T) {

	a := []string{"a", "b", "d", "e", "f"}
	b := []string{"b", "c", "d", "f", "j"}
	o := []string{"b", "d", "f"}

	out := IntersectionOf(a, b)

	if len(o) != len(out) {
		t.Errorf("Output list has incorrect length: %d", len(out))
		return
	}

	for i := range out {
		if out[i] != o[i] {
			t.Errorf("out[%d]: %q did not match expected[%d]: %q", i, out[i], i, o[i])
		}
	}
}
