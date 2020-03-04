package lioss

import (
	"fmt"
	"io"
)

/*
NGram is an implementation type of Algorithm.
*/
type NGram struct {
	nValue int
}

/*
NewNGram creates an instance of n-gram.
*/
func NewNGram(n int) *NGram {
	ngram := new(NGram)
	ngram.nValue = n
	return ngram
}

/*
Prepare of NGram do nothing.
*/
func (ngram *NGram) Prepare(db *Database) error {
	return nil
}

func (ngram *NGram) String() string {
	return fmt.Sprintf("%dgram", ngram.nValue)
}

/*
Parse parses given data and create an instance of License by n-gram.
*/
func (ngram *NGram) Parse(reader io.Reader, licenseName string) (*License, error) {
	result, err := readFully(reader)
	if err != nil {
		return nil, err
	}
	return ngram.buildNGram(result, licenseName)
}

func (ngram *NGram) buildNGram(result, licenseName string) (*License, error) {
	freq := map[string]int{}
	len := len(result) - ngram.nValue + 1
	data := []byte(result)
	for i := 0; i < len; i++ {
		ngram := string(data[i : i+ngram.nValue])
		value, ok := freq[ngram]
		if !ok {
			value = 0
		}
		freq[ngram] = value + 1
	}
	return newLicense(licenseName, freq), nil
}

/*
Compare computes similarity between given two licenses.
*/
func (ngram *NGram) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}
