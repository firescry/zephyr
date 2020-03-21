package toolbox

import (
	"testing"
)

func TestAlmostEqual(t *testing.T) {
	a, b := 1.000000001, 1.000000002
	if !AlmostEqual(a, b) {
		t.Errorf("%.9f and %.9f are almostEqual", a, b)
	}
	c, d := 1.000000001, 1.000000003
	if AlmostEqual(c, d) {
		t.Errorf("%.9f and %.9f are not almostEqual", c, d)
	}
}

func TestCompSlices(t *testing.T) {
	a, b := []float64{0.0}, []float64{0.0, 0.0}
	if CompareSlices(a, b) {
		t.Errorf("%.1f and %.1f are not the same", a, b)
	}
	c, d := []float64{0.0}, []float64{1.0}
	if CompareSlices(c, d) {
		t.Errorf("%.1f and %.1f are not the same", c, d)
	}
	e, f := []float64{0.0, 1.0}, []float64{0.0, 1.0}
	if !CompareSlices(e, f) {
		t.Errorf("%.1f and %.1f are the same", e, f)
	}
}
