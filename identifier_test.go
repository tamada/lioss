package lioss

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Example_Identifier() {
	db, err := LoadDatabase(OSI_APPROVED_DATABASE)
	if err != nil {
		return
	}
	identifier, err := NewIdentifier("9gram", 0.75, db)
	if err != nil {
		return
	}
	project, err := NewProject("testdata/project2")
	if err != nil {
		return
	}
	resultMap, err := identifier.Identify(project)
	if err != nil {
		return
	}
	for k, results := range resultMap {
		fmt.Println(k)
		for _, result := range results {
			fmt.Printf("\t%s\n", result.String())
		}
	}
	// Output:
	// license.txt
	// 	GPL-3.0-only (0.980265)
	//	GPL-3.0-or-later (0.980265)
	//	AGPL-3.0-only (0.965360)
	//	AGPL-3.0-or-later (0.965360)
}

func TestNewIdentifier(t *testing.T) {
	testdata := []struct {
		algorithm   string
		threshold   float64
		successFlag bool
	}{
		{"5gram", 0.75, true},
		{"unknown", 0.75, false},
	}
	db, _ := ReadDatabase("testdata/test.liossdb")
	for _, td := range testdata {
		_, err := NewIdentifier(td.algorithm, td.threshold, db)
		if (err == nil) != td.successFlag {
			t.Errorf("NewIdentifier(%s, %f) did not match, wont %v", td.algorithm, td.threshold, td.successFlag)
		}
	}
}

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
	db, _ := ReadDatabase("testdata/test.liossdb")
	for _, td := range testdata {
		identifier, _ := NewIdentifier(td.algorithm, td.threshold, db)
		license, _ := identifier.readLicense(createLicenseFile(td.givePath))
		results, err := identifier.identify(license)
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
	return &basicLicenseFile{id: name, reader: file}
}
