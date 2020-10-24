package lioss

import (
	"fmt"
	"io"
)

/*
nGram is an implementation type of Algorithm.
*/
type nGram struct {
	nValue int
}

/*
newNGram creates an instance of n-gram.
*/
func newNGram(n int) *nGram {
	ngram := new(nGram)
	ngram.nValue = n
	return ngram
}

/*
Prepare of NGram do nothing.
*/
func (ngram *nGram) Prepare(db *Database) error {
	return nil
}

func (ngram *nGram) String() string {
	return fmt.Sprintf("%dgram", ngram.nValue)
}

/*
Parse parses given data and create an instance of License by n-gram.
*/
func (ngram *nGram) Parse(reader io.Reader, licenseName string) (*License, error) {
	result, err := readFully(reader)
	if err != nil {
		return nil, err
	}
	return ngram.buildNGram(result, licenseName)
}

func ngramFrequency(freq map[string]int, ngram string) int {
	value, ok := freq[ngram]
	if !ok {
		value = 0
	}
	return value
}

func (ngram *nGram) buildNGram(result, licenseName string) (*License, error) {
	freq := map[string]int{}
	len := len(result) - ngram.nValue + 1
	data := []byte(result)
	for i := 0; i < len; i++ {
		ngramKey := string(data[i : i+ngram.nValue])
		value := ngramFrequency(freq, ngramKey)
		freq[ngramKey] = value + 1
	}
	return newLicense(licenseName, freq), nil
}

/*
Compare computes similarity between given two licenses.
*/
func (ngram *nGram) Compare(license1, license2 *License) float64 {
	return license1.similarity(license2)
}
