package lioss

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestMasterLicenses(t *testing.T) {
	identifier := &Identifier{}
	identifier.BuildMasterLicenses("data")
	if len(identifier.master.projects) != 20 {
		t.Errorf("master data wont 20, but %d", len(identifier.master.projects))
	}
}

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
		reader, _ := project.Open()
		resultString, err := readFully(reader)
		if err != nil {
			for _, item := range td.contains {
				if strings.Index(resultString, item) < 0 {
					t.Errorf("%s does not contain the phrase %s", td.projectDir, item)
				}
			}
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

func buildLicense(filePath string, algorithm Algorithm) *License {
	project := &BasicProject{baseDir: filepath.Dir(filePath), licenseFile: filePath}
	reader, _ := project.Open()
	license, _ := algorithm.Parse(reader, filepath.Base(filePath))
	return license
}

func TestCompare(t *testing.T) {
	algorithm, _ := CreateAlgorithm("9gram")
	licenses := map[string]*License{
		"wtfpl": buildLicense("data/WTFPL", algorithm), "gpl3": buildLicense("data/GPLv3.0", algorithm),
		"apache": buildLicense("data/Apache-License-2.0", algorithm), "bsd2": buildLicense("data/BSD-2-Clause", algorithm),
	}
	identifier := &Identifier{Threshold: 0.75, Algorithm: algorithm}

	pairs := []struct {
		first     string
		second    string
		threshold float64
	}{
		{"wtfpl", "wtfpl", 0.95},
		{"wtfpl", "gpl3", 0.001},
	}
	for _, pair := range pairs {
		result := identifier.Compare(licenses[pair.first], licenses[pair.second])
		// fmt.Printf("compare(%s, %s) is %f (%f)\n", pair.first, pair.second, result, pair.threshold)
		if result < pair.threshold {
			t.Errorf("compare(%s, %s) should be greater than %f, but got %f", pair.first, pair.second, pair.threshold, result)
		}
	}
}
