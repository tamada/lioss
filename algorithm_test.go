package lioss

import (
	"testing"
)

func TestCreateAlgorithm(t *testing.T) {
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
		algorithm, err := CreateAlgorithm(td.giveString)
		if (err == nil) != td.successFlag {
			t.Errorf("Invalid result in CreateAlgorithm, wont %v", td.successFlag)
		}
		if err == nil && algorithm.String() != td.giveString {
			t.Errorf("invalid algorithm name, wont %s, got %s", td.giveString, algorithm.String())
		}
	}
}
