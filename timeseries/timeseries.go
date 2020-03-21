package timeseries

import (
	"math"
)

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

func GenerateWeights(length int) []float64 {
	weights := make([]float64, length)
	a := 2.0 / (float64(length) + 1.0)
	for i := 0; i < length; i++ {
		weights[length-1-i] = math.Pow(1.0-a, float64(i))
	}
	return weights
}

func Initialize(length int, initValue float64) (result []float64) {
	result = make([]float64, length)
	for i := 0; i < length; i++ {
		result[i] = initValue
	}
	return result
}

func AddValue(timeSeries []float64, value float64) (result []float64) {
	result = timeSeries[1:]
	result = append(result, value)
	return result
}

func WeightedAverage(values, weights []float64) (result float64) {
	result = sumElements(multiplySlices(values, weights))
	result = result / sumElements(weights)
	return result
}
