package lioss

import "math"

type License struct {
	frequencies map[string]int
}

func (license *License) Of(key string) int {
	return license.frequencies[key]
}

func (license *License) Magnitude() float64 {
	sum := 0
	for _, value := range license.frequencies {
		sum += value * value
	}
	return math.Sqrt(float64(sum))
}
