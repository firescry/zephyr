package toolbox

import (
	"math"
)

const almostEqualThreshold = 1e-9

// AlmostEqual - compare floats
func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= almostEqualThreshold
}

// CompSlices - compare two slices with floats
func CompareSlices(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !AlmostEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
