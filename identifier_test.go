package lioss

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIdentifier(t *testing.T) {
	testdata := []struct {
		algorithm   string
		threshold   float64
		givePath    string
		successFlag bool
		wontCount   int
	}{
		{"5gram", 0.75, "LICENSE", true, 1},
	}
	db, _ := LoadDatabase("testdata/liossdb.json")
	for _, td := range testdata {
		identifier, _ := NewIdentifier(td.algorithm, td.threshold, db)
		license, _ := identifier.ReadLicense(createLicenseFile(td.givePath))
		results, err := identifier.Identify(license)
		if (err == nil) != td.successFlag {
			t.Errorf("the result of identify (%s, %s) did not match, wont %v", td.algorithm, td.givePath, td.successFlag)
		}
		if err == nil && len(results) != td.wontCount {
			t.Errorf("result count of (%s, %s) did not match, wont %d, got %d", td.algorithm, td.givePath, td.wontCount, len(results))
		}
	}
}

func createLicenseFile(path string) LicenseFile {
	name := filepath.Base(path)
	file, _ := os.Open(path)
	return &BasicLicenseFile{id: name, reader: file}
}
