package fancurve

import (
	"reflect"
	"testing"
)

func TestCurve(t *testing.T) {
	inv := []float64{20.0, 30.0, 40.0, 50.0, 60.0, 70.0, 80.0, 90.0, 100.0}
	exp := []float64{32.0, 32.0, 32.0, 49.0, 66.0, 83.0, 100.0, 100.0, 100.0}
	got := make([]float64, len(inv))
	for i, v := range inv {
		got[i] = Curve(v)
	}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("got %.1f; expected %.1f", got, exp)
	}
}

func BenchmarkCurve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Curve(55.5)
	}
}
