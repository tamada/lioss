package lioss

import (
	"math"
)

/*
License shows the license data for identifying.
*/
type License struct {
	Name        string         `json:"license-name"`
	Frequencies map[string]int `json:"frequencies"`
}

func newLicense(name string, data map[string]int) *License {
	return &License{Name: name, Frequencies: data}
}

func extractKeys(license1, license2 *License) map[string]int {
	keys := map[string]int{}
	for key := range license1.Frequencies {
		keys[key] = 1
	}
	for key := range license2.Frequencies {
		keys[key] = 1
	}
	return keys
}

func (license *License) total() int {
	sum := 0
	for _, count := range license.Frequencies {
		sum += count
	}
	return sum
}

/*
Similarity calculates the similarity between license and other by cosine similarity.
*/
func (license *License) Similarity(other *License) float64 {
	keys := extractKeys(license, other)
	sum := 0
	for key := range keys {
		sum += (license.Frequencies[key] * other.Frequencies[key])
	}
	return float64(sum) / (license.Magnitude() * other.Magnitude())
}

/*
Magnitude calculates the length of license.
*/
func (license *License) Magnitude() float64 {
	sum := 0
	for _, value := range license.Frequencies {
		sum += value * value
	}
	return math.Sqrt(float64(sum))
}
