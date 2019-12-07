package lioss

import (
	"strings"
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

func TestReadfully(t *testing.T) {
	testdata := []struct {
		projectDir string
		contains   []string
	}{
		{"testdata/project3", []string{"DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE", "Version 2, December 2004", "Everyone is permitted to copy and distribute"}},
	}
	for _, td := range testdata {
		project := NewBasicProject(td.projectDir)
		resultString, err := readFully(project)
		if err != nil {
			for _, item := range td.contains {
				if strings.Index(resultString, item) < 0 {
					t.Errorf("%s does not contain the phrase %s", td.projectDir, item)
				}
			}
		}
	}
}
