package lioss

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type Algorithm interface {
	Parse(reader io.Reader, licenseName string) (*License, error)
	Compare(license1, license2 *License) float64
	String() string
}

type Tfidf struct {
}

type NGram struct {
	nValue int
}

func CreateAlgorithm(name string) (Algorithm, error) {
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, "gram") {
		nString := strings.Replace(lowerName, "gram", "", -1)
		value, err := strconv.Atoi(nString)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid algorithm name, %s", name, err.Error())
		}
		return NewNGram(value), nil
	} else if lowerName == "tfidf" {
		return NewTfidf(), nil
	}
	return nil, fmt.Errorf("%s: unknown algorithm", lowerName)
}

func NewTfidf() *Tfidf {
	return new(Tfidf)
}

func (tfidf *Tfidf) String() string {
	return "tfidf"
}

func (tfidf *Tfidf) Parse(reader io.Reader, licenseName string) (*License, error) {
	return nil, nil
}

func (tfidf *Tfidf) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}

func NewNGram(n int) *NGram {
	ngram := new(NGram)
	ngram.nValue = n
	return ngram
}

func (ngram *NGram) String() string {
	return fmt.Sprintf("%dgram", ngram.nValue)
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
	return NewLicense(licenseName, freq), nil
}

func (ngram *NGram) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}
