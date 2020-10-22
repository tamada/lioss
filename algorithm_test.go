package lioss

import (
	"testing"
)

func TestCreateComparator(t *testing.T) {
	testdata := []struct {
		giveString  string
		successFlag bool
	}{
		{"5gram", true},
		{"kgram", false},
		{"hoge", false},
		{"wordfreq", true},
		{"tfidf", true},
	}
	for _, td := range testdata {
		comparator, err := NewAlgorithm(td.giveString)
		if (err == nil) != td.successFlag {
			t.Errorf("Invalid result in CreateComparator, wont %v", td.successFlag)
		}
		if err == nil && comparator.String() != td.giveString {
			t.Errorf("invalid comparator name, wont %s, got %s", td.giveString, comparator.String())
		}
	}
}
