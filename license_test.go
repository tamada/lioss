package lioss

import (
	"io"
	"os"
	"testing"
)

const delta float64 = 0.001

func readAll(path string) io.Reader {
	file, _ := os.Open(path)
	return file
}

func TestLicenseSimilarity(t *testing.T) {
	testdata := []struct {
		comparator string
		path1      string
		path2      string
		similarity float64
	}{
		{"5gram", "data/misc/WTFPL", "data/misc/WTFPL", 1.0},
		{"5gram", "data/misc/BSD-3-Clause", "data/misc/BSD-4-Clause", 0.9385},
		{"wordfreq", "data/misc/WTFPL", "data/misc/WTFPL", 1.0},
	}
	for _, td := range testdata {
		comparator, _ := CreateComparator(td.comparator)
		license1, _ := comparator.Parse(readAll(td.path1), "license1")
		license2, _ := comparator.Parse(readAll(td.path2), "license2")
		similarity := comparator.Compare(license1, license2)
		if td.similarity-delta > similarity || td.similarity+delta < similarity {
			t.Errorf("similarity between %s and %s no in the suitable range, wont %f (%f), got %f", td.path1, td.path2, td.similarity, delta, similarity)
		}
	}
}
