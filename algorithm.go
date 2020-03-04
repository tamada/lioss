package lioss

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

/*
Algorithm shows an algorithm for identifying the license.
*/
type Algorithm interface {
	Prepare(db *Database) error
	Parse(reader io.Reader, licenseName string) (*License, error)
	Compare(license1, license2 *License) float64
	String() string
}

/*
CreateAlgorithm create an instance of Algorithm.
Available values are [1-9]gram, and tfidf.
*/
func CreateAlgorithm(name string) (Algorithm, error) {
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, "gram") {
		nString := strings.Replace(lowerName, "gram", "", -1)
		value, err := strconv.Atoi(nString)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid algorithm name, %s", name, err.Error())
		}
		return NewNGram(value), nil
	} else if lowerName == "wordfreq" {
		return NewWordFreq(), nil
	} else if lowerName == "tfidf" {
		return NewTfidf(), nil
	}
	return nil, fmt.Errorf("%s: unknown algorithm", lowerName)
}

func normalize(dataArray []byte) string {
	data := strings.ReplaceAll(string(dataArray), "\r", " ")
	data = strings.ReplaceAll(data, "\n", " ")
	data = strings.ReplaceAll(data, "\t", " ")
	for strings.Index(data, "  ") >= 0 {
		data = strings.ReplaceAll(data, "  ", " ")
	}
	return strings.TrimSpace(data)
}

func readFully(reader io.Reader) (string, error) {
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return normalize(result), nil
}
