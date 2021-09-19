package fancurve

import (
	"reflect"
	"testing"
)

func TestNewCurve(t *testing.T) {
	curve := NewCurve(40.0, 80.0, 32.0, 100.0)
	if curve.a != 1.7 || curve.b != -36.0 {
		t.Errorf("got a=%.1f, b=%.1f; expected a=%.1f, b=%.1f", curve.a, curve.b, 1.7, -36.0)
	}
}

func TestGetPwmPercent(t *testing.T) {
	curve := NewCurve(40.0, 80.0, 32.0, 100.0)
	inv := []float64{20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0}
	exp := []float64{32.0, 32.0, 32.0, 49.0, 66.0, 83.0, 100.0, 100.0, 100.0}
	got := make([]float64, len(inv))
	for i, v := range inv {
		got[i] = curve.GetPwmPercent(v)
	}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("got %.1f; expected %.1f", got, exp)
	}
}

func BenchmarkGetPwmPercent(b *testing.B) {
	curve := NewCurve(40.0, 80.0, 32.0, 100.0)
	for i := 0; i < b.N; i++ {
		curve.GetPwmPercent(55.5)
	}
}
