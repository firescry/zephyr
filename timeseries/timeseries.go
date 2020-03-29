package timeseries

import (
	"math"
)

type TimeSeries struct {
	samples []float64
	weights []float64
}

func sumElements(a []float64) (sum float64) {
	for _, elementValue := range a {
		sum += elementValue
	}
	return sum
}

func multiplySlices(a, b []float64) (result []float64) {
	result = make([]float64, len(a))
	for i := range a {
		result[i] = a[i] * b[i]
	}
	return result
}

func InitTimeSeries(length int, initValue float64) *TimeSeries {
	timeSeries := TimeSeries{}
	timeSeries.initSamples(length, initValue)
	timeSeries.initWeights(length)
	return &timeSeries
}

func (timeSeries *TimeSeries) initSamples(length int, initValue float64) {
	timeSeries.samples = make([]float64, length)
	for i := 0; i < length; i++ {
		timeSeries.samples[i] = initValue
	}
}

func (timeSeries *TimeSeries) initWeights(length int) {
	timeSeries.weights = make([]float64, length)
	a := 2.0 / (float64(length) + 1.0)
	for i := 0; i < length; i++ {
		timeSeries.weights[length-1-i] = math.Pow(1.0-a, float64(i))
	}
}

func (timeSeries *TimeSeries) AddSample(sample float64) {
	timeSeries.samples = timeSeries.samples[1:]
	timeSeries.samples = append(timeSeries.samples, sample)
}

func (timeSeries *TimeSeries) WeightedAverage() (result float64) {
	result = sumElements(multiplySlices(timeSeries.samples, timeSeries.weights))
	result = result / sumElements(timeSeries.weights)
	return result
}
