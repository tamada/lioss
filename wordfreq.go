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

func (ngram *WordFreq) String() string {
	return "wordfreq"
}

func preprocess(str string) string {
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

func buildLicense(licenseName string, words []string) (*License, error) {
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
func (ngram *WordFreq) Parse(reader io.Reader, licenseName string) (*License, error) {
	result, err := readFully(reader)
	if err != nil {
		return nil, err
	}
	result = preprocess(result)
	return buildLicense(licenseName, strings.Split(result, " "))
}

/*
Compare computes similarity between given two licenses.
*/
func (ngram *WordFreq) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}
