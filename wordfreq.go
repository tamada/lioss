package lioss

import (
	"io"
	"strings"
)

/*
NGram is an implementation type of Algorithm.
*/
type WordFreq struct {
}

/*
NewNGram creates an instance of n-gram.
*/
func NewWordFreq() *WordFreq {
	return new(WordFreq)
}

func (wfreq *WordFreq) String() string {
	return "wordfreq"
}

/*
Prepare of WordFreq do nothing.
*/
func (wfreq *WordFreq) Prepare(db *Database) error {
	return nil
}

func preprocessForWordFreq(str string) string {
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, ";", "")
	str = strings.ReplaceAll(str, "!", "")
	str = strings.ReplaceAll(str, "?", "")
	str = strings.ReplaceAll(str, "`", "")
	str = strings.ReplaceAll(str, "<", "")
	str = strings.ReplaceAll(str, ">", "")
	str = strings.ReplaceAll(str, "(", "")
	str = strings.ReplaceAll(str, ")", "")
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "\"", "")
	return strings.ToLower(str)
}

func BuildWordFreqLicense(licenseName string, document string) (*License, error) {
	document = preprocessForWordFreq(document)
	words := strings.Split(document, " ")
	freq := map[string]int{}
	for _, word := range words {
		count, ok := freq[word]
		if !ok {
			count = 0
		}
		freq[word] = count + 1
	}
	return newLicense(licenseName, freq), nil
}

/*
Parse parses given data and create an instance of License by n-gram.
*/
func (wfreq *WordFreq) Parse(reader io.Reader, licenseName string) (*License, error) {
	result, err := readFully(reader)
	if err != nil {
		return nil, err
	}
	return BuildWordFreqLicense(licenseName, result)
}

/*
Compare computes similarity between given two licenses.
*/
func (wfreq *WordFreq) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}
