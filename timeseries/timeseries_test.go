package timeseries

import (
	"reflect"
	"testing"

	"github.com/firescry/zephyr/toolbox"
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

func TestInitTimeSeries(t *testing.T) {
	exp := TimeSeries{
		samples: []float64{1.0, 1.0, 1.0},
		weights: []float64{0.25, 0.5, 1.0},
	}
	length := 3
	init := 1.0
	got := *InitTimeSeries(length, init)
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("InitTimeSeries(%d, %.2f) = %.2f; expected %.2f", length, init, got, exp)
	}
}

func TestInitSamples(t *testing.T) {
	timeSeries := TimeSeries{}
	length := 5
	initValue := 5.0
	timeSeries.initSamples(length, initValue)
	exp := []float64{5.0, 5.0, 5.0, 5.0, 5.0}
	if !toolbox.CompareSlices(timeSeries.samples, exp) {
		t.Errorf("timeSeries.initSamples(%d, %f) = %f; expected %f", length, initValue, timeSeries.samples, exp)
	}
}

func TestInitWeights(t *testing.T) {
	timeSeries := TimeSeries{}
	length := 5
	timeSeries.initWeights(length)
	exp := []float64{16.0 / 81.0, 8.0 / 27.0, 4.0 / 9.0, 2.0 / 3.0, 1.0}
	if !toolbox.CompareSlices(timeSeries.weights, exp) {
		t.Errorf("timeSeries.initWeights(%d) = %f; expected %f", length, timeSeries.weights, exp)
	}
}

func TestAddSample(t *testing.T) {
	inv := TimeSeries{
		samples: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		weights: []float64{},
	}
	val := 6.0
	exp := []float64{2.0, 3.0, 4.0, 5.0, 6.0}
	inv.AddSample(val)
	if !toolbox.CompareSlices(inv.samples, exp) {
		t.Errorf("inv.AddSample(%.2f) = %.2f; expected %.2f", val, inv.samples, exp)
	}
}

func BenchmarkAddSample(b *testing.B) {
	inv := InitTimeSeries(benchSliceLength, 5.4321)
	val := 9.8765
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv.AddSample(val)
	}
}

func TestWeightedAverage(t *testing.T) {
	inv := TimeSeries{
		samples: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
		weights: []float64{0.5, 1.0, 1.5, 3.0, 4.0},
	}
	exp := 3.9
	got := inv.WeightedAverage()
	if !toolbox.AlmostEqual(got, exp) {
		t.Errorf("inv.WeightedAverage() = %.2f; expected %.2f", got, exp)
	}
}

func BenchmarkWeightedAverage(b *testing.B) {
	inv := InitTimeSeries(benchSliceLength, 5.4321)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inv.WeightedAverage()
	}
}
