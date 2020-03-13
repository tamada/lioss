package lioss

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/tamada/lioss/lib"
)

/*
AvailableAlgorithms contains the names of available algorithm for comparing licenses.
*/
var AvailableAlgorithms = []string{"1gram", "2gram", "3gram", "4gram", "5gram", "6gram", "7gram", "8gram", "9gram", "wordfreq", "tfidf"}

/*
Comparator shows an algorithm for identifying the license.
*/
type Comparator interface {
	Prepare(db *Database) error
	Parse(reader io.Reader, licenseName string) (*License, error)
	Compare(license1, license2 *License) float64
	String() string
}

func createNGramComparator(name string) (Comparator, error) {
	lowerName := strings.ToLower(name)
	nString := strings.Replace(lowerName, "gram", "", -1)
	value, err := strconv.Atoi(nString)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid algorithm name, %s", name, err.Error())
	}
	return NewNGram(value), nil
}

/*
CreateComparator create an instance of Algorithm.
Available values are [1-9]gram, and tfidf.
*/
func CreateComparator(name string) (Comparator, error) {
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, "gram") {
		return createNGramComparator(name)
	} else if lowerName == "wordfreq" {
		return NewWordFreq(), nil
	} else if lowerName == "tfidf" {
		return NewTfidf(), nil
	}
	return nil, fmt.Errorf("%s: unknown algorithm", lowerName)
}

func readFully(reader io.Reader) (string, error) {
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return lib.Normalize(result), nil
}
