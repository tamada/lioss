package lioss

import (
	"bytes"
	"testing"
)

func TestDatabaseTypeString(t *testing.T) {
	testdata := []struct {
		giveType   DatabaseType
		wontString string
	}{
		{OSI_APPROVED_DATABASE, "OSI_APPROVED_DATABASE"},
		{OSI_DEPRECATED_DATABASE, "OSI_DEPRECATED_DATABASE"},
		{NONE_OSI_APPROVED_DATABASE, "NONE_OSI_APPROVED_DATABASE"},
		{DEPRECATED_DATABASE, "DEPRECATED_DATABASE"},
		{WHOLE_DATABASE, "WHOLE_DATABASE"},
		{-1, "UNKNOWN"},
	}
	for _, td := range testdata {
		if td.giveType.String() != td.wontString {
			t.Errorf("%s.String() did not match, wont %s, got %s", td.giveType, td.wontString, td.giveType.String())
		}
	}
}

func TestExtension(t *testing.T) {
	testdata := []struct {
		giveString string
		wontString string
	}{
		{"file.liossdb", "file.liossdb"},
		{"file.liossgz", "file.liossgz"},
		{"file.liossdb.gz", "file.liossgz"},
		{"file.json.gz", "file.json.liossgz"},
		{"file", "file.liossdb"},
		{"file.json", "file.liossdb"},
	}
	for _, td := range testdata {
		gotString := destination(td.giveString)
		if gotString != td.wontString {
			t.Errorf("result of normalizeDestination(%s) did not match, wont %s, got %s", td.giveString, td.wontString, gotString)
		}
	}
}

func TestLoadDatabase(t *testing.T) {
	testdata := []struct {
		dbTypes           DatabaseType
		wontLicenseSize   int
		wontAlgorithmSize int
	}{
		{NONE_OSI_APPROVED_DATABASE, 292, 11},
		{OSI_APPROVED_DATABASE, 120, 11},
		{OSI_DEPRECATED_DATABASE, 0, 0},
		{DEPRECATED_DATABASE, 18, 11},
		{WHOLE_DATABASE, 430, 11},
	}
	for _, td := range testdata {
		db, _ := LoadDatabase(td.dbTypes)
		if db.AlgorithmCount() != td.wontAlgorithmSize {
			t.Errorf("%s: algorithm count did not match: wont %d, got %d", td.dbTypes, td.wontAlgorithmSize, db.AlgorithmCount())
		}
		if db.LicenseCount() != td.wontLicenseSize {
			t.Errorf("%s: license count did not match, wont %d, got %d", td.dbTypes, td.wontLicenseSize, db.LicenseCount())
		}
	}
}

func TestReadDatabase(t *testing.T) {
	testdata := []struct {
		givePath      string
		wontDataSize  int
		wontEntrySize int
	}{
		{"testdata/test.liossdb", 11, 23},
		{"testdata/test.liossgz", 11, 23},
	}
	for _, td := range testdata {
		db, err := ReadDatabase(td.givePath)
		if err != nil {
			t.Errorf("error on LoadDatabase of %s: %s", td.givePath, err.Error())
		}
		if len(db.Data) != td.wontDataSize {
			t.Errorf("size of loaded database did not match, wont %d, got %d", td.wontDataSize, len(db.Data))
		}
		entries := db.Entries("5gram")
		if len(entries) != td.wontEntrySize {
			t.Errorf("size of 5gram entries did not match, wont %d, got %d", td.wontEntrySize, len(entries))
		}
	}
}

func TestLoadDatabaseFail(t *testing.T) {
	_, err := ReadDatabase("testdata/notexistdb.liossdb")
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
	db2, _ := Read(buffer, "memory")
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
	db, _ := ReadDatabase("testdata/test.liossdb")
	for _, td := range testdata {
		containsFlag := db.Contains(td.algorithmName, td.licenseName)
		if containsFlag != td.wontFlag {
			t.Errorf("Contains(%s, %s) did not match, wont %v, got %v", td.algorithmName, td.licenseName, td.wontFlag, containsFlag)
		}
	}
}

func TestPutLicenseToNewAlgorithm(t *testing.T) {
	db, _ := ReadDatabase("testdata/test.liossdb")
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
	db, _ := ReadDatabase("testdata/test.liossdb")
	license := newLicense("NYSL", map[string]int{"nirunari": 1, "yakunari": 1, "sukinishiro": 1, "license": 1})
	db.Put("5gram", license)
	if !db.Contains("5gram", "NYSL") {
		t.Errorf("put license did not found (5gram/NYSL).")
	}
	if db.Contains("2gram", "NYSL") {
		t.Errorf("put license found (2gram/NYSL).")
	}
	items := db.Entries("5gram")
	if len(items) != 24 {
		t.Errorf("size of license did not match (unknown), wont 24, got %d", len(items))
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
	db, _ := ReadDatabase("testdata/test.liossdb")
	license := newLicense("GPLv3.0", map[string]int{"nirunari": 1, "yakunari": 1, "sukinishiro": 1, "license": 1})
	db.Put("5gram", license)
	if !db.Contains("5gram", "GPLv3.0") {
		t.Errorf("put license did not found (5gram/GPLv3.0).")
	}
	items := db.Entries("5gram")
	if len(items) != 23 {
		t.Errorf("size of license did not match (unknown), wont 23, got %d", len(items))
	}
	item := db.Entry("5gram", "GPLv3.0")
	if item == nil || len(item.Frequencies) != 4 {
		t.Errorf("size of map did not match, wont 4, got %d", len(item.Frequencies))
	}
}
