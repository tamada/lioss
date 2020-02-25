package lioss

import (
	"testing"
)

func TestNormalizeReading(t *testing.T) {
	testdata := []struct {
		giveString []byte
		wontString string
	}{
		{[]byte("a  b  c\n  d  e\r\nf g \th \r i j"), "a b c d e f g h i j"},
	}
	for _, td := range testdata {
		result := normalize(td.giveString)
		if result != td.wontString {
			t.Errorf("normalize(%s), wont %s, but got %s", string(td.giveString), td.wontString, result)
		}
	}
}
func TestBuildAlgorithm(t *testing.T) {
	testdata := []struct {
		giveName   string
		errorFlag  bool
		typeString string
	}{
		{"9gram", false, "9gram"},
		{"tfidf", false, "tfidf"},
		{"unknown", true, ""},
	}
	for _, td := range testdata {
		algorithm, err := CreateAlgorithm(td.giveName)
		if (err != nil) != td.errorFlag {
			t.Errorf("errorFlag of BuildAlgorithm(%s) should be %v, but %v (%v)", td.giveName, td.errorFlag, !td.errorFlag, err)
		}
		if err == nil {
			if algorithm.String() != td.typeString {
				t.Errorf("BuildAlgorithm(%s) created type did not match, wont %s, got %s", td.giveName, td.typeString, algorithm.String())
			}
		}
	}
}
