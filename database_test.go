package lioss

import (
	"bytes"
	"testing"
)

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

func TestLoadDatabaseFail(t *testing.T) {
	_, err := LoadDatabase("testdata/notexistdb.json")
	if err == nil {
		t.Errorf("successfully load not exist database .")
	}
}

func TestWriteLoad(t *testing.T) {
	db := NewDatabase()
	license := newLicense("hoge", map[string]int{"a": 1, "b": 2, "c": 3, "d": 4})
	db.Put("1gram", license)
	buffer := bytes.NewBuffer([]byte{})
	db.Write(buffer)
	db2, _ := Load(buffer)
	if len(db2.Data) != 1 {
		t.Errorf("db write/load error, wont len(db2.Data) == %d, got %d", 1, len(db2.Data))
	}
	if !db2.Contains("1gram", "hoge") {
		t.Errorf(`db write/load error wont db2.Contains("%s", "%s")=true, got false`, "1gram", "hoge")
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
	if entry := db.Entry("5gram", "NotExistLicense"); entry != nil {
		t.Errorf(`db.Entry("%s", "%s") returns some instance: %v`, "5gram", "NotExistLicense", entry)
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
