package lioss

import (
	"strings"
	"testing"
)

func TestWordFreq(t *testing.T) {
	testdata := []struct {
		givenData string
		results   map[string]int
	}{
		{givenData: `today is fine.
I am fine, too!`, results: map[string]int{"today": 1, "is": 1, "fine": 2, "i": 1, "am": 1, "too": 1}},
		{givenData: `The quick brown fox jumps over the lazy dog.
The quick onyx goblin jumps over the lazy dwarf.`, results: map[string]int{"the": 4, "quick": 2, "brown": 1, "fox": 1, "jumps": 2, "over": 2, "lazy": 2, "dog": 1, "onyx": 1, "goblin": 1, "dwarf": 1}},
	}
	for _, td := range testdata {
		wfreq := NewWordFreq()
		license, err := wfreq.Parse(strings.NewReader(td.givenData), "unknown-license")
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
