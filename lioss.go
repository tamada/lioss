package lioss

import (
	"io/ioutil"
	"strings"
)

type Algorithm interface {
	Parse(project Project) (*License, error)
}

type NGramAlgorithm struct {
	ngram int
}

func NewNGramAlgorithm(n int) *NGramAlgorithm {
	ngram := new(NGramAlgoirthm)
	ngram.ngram = n
	return ngram
}

func (algo *NGramAlgorithm) Parse(project Project) (*License, error) {
	result, err := readFully(project)
	if err != nil {
		return nil, err
	}
	return buildNGram(algo.ngram, result)
}

func buildNGram(n int, result string) (*License, error) {
	freq := map[string]int{}
	len := len(result) - n + 1
	data := []byte(result)
	for i := 0; i < len; i++ {
		ngram := string(data[i : i+n])
		value, ok := freq[ngram]
		if !ok {
			value = 0
		}
		freq[ngram] = value + 1
	}
	return &License{frequencies: freq}, nil
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

func readFully(project Project) (string, error) {
	reader, err := project.Open()
	if err != nil {
		return "", err
	}
	defer project.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	result := normalize(data)
	return result, nil
}
