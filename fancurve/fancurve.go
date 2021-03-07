package fancurve

// Curve returns desired fan speed (in percentage) for given temperature
func Curve(temp float64) float64 {
	switch {
	case temp < 40.0:
		return 32.0
	case temp > 80.0:
		return 100.0
	default:
		return 1.7*temp - 36.0
	}
}
