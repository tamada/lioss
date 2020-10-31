package lioss

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func readAllLicenses(tfidf Algorithm, fromDir string) []*License {
	licenses := []*License{}
	files, _ := ioutil.ReadDir(fromDir)
	for _, file := range files {
		path := filepath.Join(fromDir, file.Name())
		reader, _ := os.Open(path)
		defer reader.Close()
		license, _ := tfidf.Parse(reader, file.Name())
		licenses = append(licenses, license)
	}
	return licenses
}

func CreateDB(tfidf Algorithm, fromDir, toPath string) *Database {
	licenses := readAllLicenses(tfidf, fromDir)
	db := NewDatabase()
	db.Data = map[string][]*License{}
	db.Data["tfidf"] = licenses
	writer, _ := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY, 0644)
	defer writer.Close()
	db.Write(writer)
	return db
}

func createLicenseFromFile(tfidf Algorithm, path string) *License {
	reader, _ := os.Open(path)
	defer reader.Close()
	license, _ := tfidf.Parse(reader, "")
	return license
}

func TestCompare(t *testing.T) {
	tfidf := newTfidf()
	db := CreateDB(tfidf, "data/misc", "tfidf.json")
	defer os.Remove("tfidf.json")
	tfidf.Prepare(db)
	testdata := []struct {
		license        string
		path           string
		wontSimilarity float64
	}{
		{"MIT", "data/misc/MIT", 0.99},
		{"WTFPL", "data/misc/WTFPL", 0.95},
		{"WTFPL", "testdata/project1/LICENSE", 0.9},
	}
	for _, td := range testdata {
		gotSimilarity := tfidf.Compare(&License{Name: td.license}, createLicenseFromFile(tfidf, td.path))
		if gotSimilarity < td.wontSimilarity {
			t.Errorf("tfidf.Compare(%s, %s) is less than %f, was %f", td.license, td.path, td.wontSimilarity, gotSimilarity)
		}
	}
}

func TestStoreTfidfData(t *testing.T) {
	testdata := []struct {
		givenData string
		results   map[string]int
	}{
		{givenData: `today is fine.
I am fine, too!`, results: map[string]int{"today": 1, "is": 1, "fine": 2, "i": 1, "am": 1, "too": 1}},
		{givenData: `The quick brown fox jumps over the lazy dog.
The quick onyx goblin jumps over the lazy dwarf.`, results: map[string]int{"the": 4, "quick": 2, "brown": 1, "fox": 1, "jumps": 2, "over": 2, "lazy": 2, "dog": 1, "onyx": 1, "goblin": 1, "dwarf": 1}},
	}
	for _, td := range testdata {
		tfidf := newTfidf()
		license, err := tfidf.Parse(strings.NewReader(td.givenData), "unknown-license")
		if err != nil {
			t.Errorf("parse failed: %s", err.Error())
		}
		if len(license.Frequencies) != len(td.results) {
			t.Errorf("map size did not match, wont %d, got %d (%v", len(td.results), len(license.Frequencies), license)
		}
		for key, wontValue := range td.results {
			gotValue := license.Frequencies[key]
			if wontValue != gotValue {
				t.Errorf("value did not match, wont: ngram[%s]: %d, got ngram[%s]: %d", key, wontValue, key, gotValue)
			}
		}
	}
}
