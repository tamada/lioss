package lioss

import "testing"

func TestLoadDatabase(t *testing.T) {
	db, err := LoadDatabase("testdata/liossdb.json")
	if err != nil {
		t.Errorf("error on LoadDatabase: %s", err.Error())
	}
	if len(db.Data) != 9 {
		t.Errorf("size of loaded database did not match, wont %d, got %d", 9, len(db.Data))
	}
	entries := db.Entries("5gram")
	if len(entries) != 20 {
		t.Errorf("size of 5gram entries did not match, wont %d, got %d", 20, len(entries))
	}
}

func TestContains(t *testing.T) {
	testdata := []struct {
		algorithmName string
		licenseName   string
		wontFlag      bool
	}{
		{"5gram", "GPLv3.0", true},
		{"10gram", "GPLv3.0", false},
		{"2gram", "GPL", false},
	}
	db, _ := LoadDatabase("testdata/liossdb.json")
	for _, td := range testdata {
		containsFlag := db.Contains(td.algorithmName, td.licenseName)
		if containsFlag != td.wontFlag {
			t.Errorf("Contains(%s, %s) did not match, wont %v, got %v", td.algorithmName, td.licenseName, td.wontFlag, containsFlag)
		}
	}
}

func TestPutLicenseToNewAlgorithm(t *testing.T) {
	db, _ := LoadDatabase("testdata/liossdb.json")
	license := newLicense("NYSL", map[string]int{"nirunari": 1, "yakunari": 1, "sukinishiro": 1, "license": 1})
	db.Put("unknown", license)
	if !db.Contains("unknown", "NYSL") {
		t.Errorf("put license did not found (unknown/NYSL).")
	}
	items := db.Entries("unknown")
	if len(items) != 1 {
		t.Errorf("size of license did not match (unknown), wont 1, got %d", len(items))
	}
	if len(items[0].Frequencies) != 4 {
		t.Errorf("size of map did not match, wont 4, got %d", len(items[0].Frequencies))
	}
}

func TestPutLicenseTo5GramWithNotContainedLicense(t *testing.T) {
	db, _ := LoadDatabase("testdata/liossdb.json")
	license := newLicense("NYSL", map[string]int{"nirunari": 1, "yakunari": 1, "sukinishiro": 1, "license": 1})
	db.Put("5gram", license)
	if !db.Contains("5gram", "NYSL") {
		t.Errorf("put license did not found (5gram/NYSL).")
	}
	if db.Contains("2gram", "NYSL") {
		t.Errorf("put license found (2gram/NYSL).")
	}
	items := db.Entries("5gram")
	if len(items) != 21 {
		t.Errorf("size of license did not match (unknown), wont 21, got %d", len(items))
	}
	item := db.Entry("5gram", "NYSL")
	if item == nil || len(item.Frequencies) != 4 {
		t.Errorf("size of map did not match, wont 4, got %d", len(item.Frequencies))
	}
}

func TestReplaceLicenseData(t *testing.T) {
	db, _ := LoadDatabase("testdata/liossdb.json")
	license := newLicense("GPLv3.0", map[string]int{"nirunari": 1, "yakunari": 1, "sukinishiro": 1, "license": 1})
	db.Put("5gram", license)
	if !db.Contains("5gram", "GPLv3.0") {
		t.Errorf("put license did not found (5gram/GPLv3.0).")
	}
	items := db.Entries("5gram")
	if len(items) != 20 {
		t.Errorf("size of license did not match (unknown), wont 20, got %d", len(items))
	}
	item := db.Entry("5gram", "GPLv3.0")
	if item == nil || len(item.Frequencies) != 4 {
		t.Errorf("size of map did not match, wont 4, got %d", len(item.Frequencies))
	}
}
