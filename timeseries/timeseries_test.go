package timeseries

import (
	"github.com/firescry/zephyr/toolbox"
	"testing"
)

const benchSliceLength = 30

func TestSumElements(t *testing.T) {
	inv := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	exp := 15.0
	got := sumElements(inv)
	if got != exp {
		t.Errorf("sumElements(%f) = %f; expected %f", inv, got, exp)
	}
}

func BenchmarkSumElements(b *testing.B) {
	inv := make([]float64, benchSliceLength)
	for i := 0; i < benchSliceLength; i++ {
		inv[i] = float64(i) * 5.4321
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sumElements(inv)
	}
}

func TestMultiplySlices(t *testing.T) {
	inv := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	exp := []float64{1.0, 4.0, 9.0, 16.0, 25.0}
	got := multiplySlices(inv, inv)
	if !toolbox.CompareSlices(got, exp) {
		t.Errorf("multiplySlices(%f, %f) = %f; expected %f", inv, inv, got, exp)
	}
}

func BenchmarkMultiplySlices(b *testing.B) {
	inv := make([]float64, benchSliceLength)
	for i := 0; i < benchSliceLength; i++ {
		inv[i] = float64(i) * 5.4321
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		multiplySlices(inv, inv)
	}
}

func TestGenerateWeights(t *testing.T) {
	exp := []float64{16.0 / 81.0, 8.0 / 27.0, 4.0 / 9.0, 2.0 / 3.0, 1.0}
	got := GenerateWeights(5)
	if !toolbox.CompareSlices(got, exp) {
		t.Errorf("GenerateWeights(5) = %f; expected %f", got, exp)
	}
}

func TestInitialize(t *testing.T) {
	length := 5
	initValue := 5.0
	exp := []float64{5.0, 5.0, 5.0, 5.0, 5.0}
	got := Initialize(length, initValue)
	if !toolbox.CompareSlices(got, exp) {
		t.Errorf("Initialize(%d, %f) = %f; expected %f", length, initValue, got, exp)
	}
}

func TestAddValue(t *testing.T) {
	inv := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	val := 6.0
	exp := []float64{2.0, 3.0, 4.0, 5.0, 6.0}
	got := AddValue(inv, val)
	if !toolbox.CompareSlices(got, exp) {
		t.Errorf("AddValue(%f, %f) = %f; expected %f", inv, val, got, exp)
	}
}

func BenchmarkAddValue(b *testing.B) {
	inv := make([]float64, benchSliceLength)
	for i := 0; i < benchSliceLength; i++ {
		inv[i] = float64(i) * 5.4321
	}
	val := 9.8765
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddValue(inv, val)
	}
}

func TestWeightedAverage(t *testing.T) {
	inv := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	wgt := []float64{0.5, 1.0, 1.5, 3.0, 4.0}
	exp := 3.9
	got := WeightedAverage(inv, wgt)
	if !toolbox.AlmostEqual(got, exp) {
		t.Errorf("WeightedAverage(%f, %f) = %f; expected %f", inv, wgt, got, exp)
	}
}

func BenchmarkWeightedAverage(b *testing.B) {
	inv := make([]float64, benchSliceLength)
	wgt := make([]float64, benchSliceLength)
	for i := 0; i < benchSliceLength; i++ {
		inv[i] = float64(i) * 5.4321
		wgt[i] = float64(i) * 9.8765
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WeightedAverage(inv, wgt)
	}
}
