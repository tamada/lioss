package lioss

import (
	"math"
	"time"
)

type License struct {
	LicenseName string
	frequencies map[string]int
	LoadDate    time.Time
}

func NewLicense(name string, data map[string]int) *License {
	return &License{LicenseName: name, frequencies: data, LoadDate: time.Now()}
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
