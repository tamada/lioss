package lioss

import (
	"io"
	"strings"
)

/*
wordFreq is an implementation type of Algorithm.
*/
type wordFreq struct {
}

/*
newWordFreq creates an instance of wordfreq.
*/
func newWordFreq() *wordFreq {
	return new(wordFreq)
}

func (wfreq *wordFreq) String() string {
	return "wordfreq"
}

/*
Prepare of WordFreq do nothing.
*/
func (wfreq *wordFreq) Prepare(db *Database) error {
	return nil
}

func preprocessForWordFreq(str string) string {
	replacer := []string{
		".", ",", ";", "!", "?", "`", "<", ">", "(", ")", "'", "\"",
	}
	for _, item := range replacer {
		str = strings.ReplaceAll(str, item, "")
	}
	return strings.ToLower(str)
}

/*
BuildWordFreqLicense creates an instance of License by wordfreq algorithm.
*/
func buildWordFreqLicense(licenseName string, document string) (*License, error) {
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
func (wfreq *wordFreq) Parse(reader io.Reader, licenseName string) (*License, error) {
	result, err := readFully(reader)
	if err != nil {
		return nil, err
	}
	return buildWordFreqLicense(licenseName, result)
}

/*
Compare computes similarity between given two licenses.
*/
func (wfreq *wordFreq) Compare(license1, license2 *License) float64 {
	return license1.similarity(license2)
}
