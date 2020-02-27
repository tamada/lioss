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
		algorithm  string
		path1      string
		path2      string
		similarity float64
	}{
		{"5gram", "data/WTFPL", "data/WTFPL", 1.0},
	}
	for _, td := range testdata {
		algorithm, _ := CreateAlgorithm(td.algorithm)
		license1, _ := algorithm.Parse(readAll(td.path1), "license1")
		license2, _ := algorithm.Parse(readAll(td.path2), "license2")
		similarity := license1.Similarity(license2)
		if td.similarity-delta > similarity || td.similarity+delta < similarity {
			t.Errorf("similarity between %s and %s no in the suitable range, wont %f (%f), got %f", td.path1, td.path2, td.similarity, delta, similarity)
		}
	}
}
