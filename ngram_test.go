package lioss

import (
	"strings"
	"testing"
)

func TestNGram(t *testing.T) {
	testdata := []struct {
		givenData string
		nValue    int
		results   map[string]int
	}{
		{givenData: "abcd", nValue: 1, results: map[string]int{"a": 1, "b": 1, "c": 1, "d": 1}},
		{givenData: "abracadabra", nValue: 3, results: map[string]int{"abr": 2, "bra": 2, "rac": 1, "aca": 1, "cad": 1, "ada": 1, "dab": 1}},
	}
	for _, td := range testdata {
		ngram := NewNGram(td.nValue)
		license, err := ngram.Parse(strings.NewReader(td.givenData), "unknown-license")
		if err != nil {
			t.Errorf("parse failed: %s", err.Error())
		}
		if len(license.Frequencies) != len(td.results) {
			t.Errorf("map size did not match, wont %d, got %d (%v", len(td.results), len(license.Frequencies), license)
		}
		for key, wontValue := range td.results {
			gotValue := license.Frequencies[key]
			if wontValue != gotValue {
				t.Errorf("value did not match, wont: ngram[%s]: %d, got ngram[%s]: %d", key, wontValue, key, gotValue)
			}
		}
	}
}
