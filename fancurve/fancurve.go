package fancurve

type Curve struct {
	minTemp float64
	maxTemp float64
	a       float64
	b       float64
}

func NewCurve(minTemp, maxTemp, minPwmPercent, maxPwmPercent float64) *Curve {
	a := (maxPwmPercent - minPwmPercent) / (maxTemp - minTemp)
	b := maxPwmPercent - a*maxTemp
	curve := Curve{
		minTemp: minTemp,
		maxTemp: maxTemp,
		a:       a,
		b:       b,
	}
	return &curve
}

func (c *Curve) GetPwmPercent(temp float64) float64 {
	if temp < c.minTemp {
		temp = c.minTemp
	} else if temp > c.maxTemp {
		temp = c.maxTemp
	}
	return c.a*temp + c.b
}
