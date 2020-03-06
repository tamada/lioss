package lioss

import (
	"io"
	"math"
)

/*
Tfidf is an implementation type of Algorithm.
*/
type Tfidf struct {
	data map[string]*document
}

type document struct {
	name  string
	words map[string]*value
}

func (doc *document) tfidf(word string) float64 {
	value, ok := doc.words[word]
	if !ok {
		return float64(0)
	}
	return value.tfidf()
}

func (doc *document) total() int {
	sum := 0
	for _, value := range doc.words {
		sum += value.count
	}
	return sum
}

func (doc *document) magnitude() float64 {
	sum := float64(0)
	for _, value := range doc.words {
		tfidf := value.tfidf()
		sum += (tfidf * tfidf)
	}
	return math.Sqrt(sum)
}

func (doc *document) contains(word string) bool {
	_, ok := doc.words[word]
	return ok
}

type value struct {
	word  string
	count int
	tf    float64
	idf   float64
}

func (val *value) tfidf() float64 {
	return val.tf * val.idf
}

/*
NewTfidf creates an instance of Tfidf.
*/
func NewTfidf() *Tfidf {
	return &Tfidf{data: map[string]*document{}}
}

func (tfidf *Tfidf) String() string {
	return "tfidf"
}

func (tfidf *Tfidf) countDocument(word string) int {
	count := 0
	for _, document := range tfidf.data {
		if document.contains(word) {
			count++
		}
	}
	return count
}

func calculateTfidf(tfidf *Tfidf, word string, count, total int) *value {
	value := &value{word: word, count: count, tf: float64(count) / float64(total)}
	value.idf = math.Log(float64(len(tfidf.data))/float64(tfidf.countDocument(word))) + float64(1)
	return value
}

func calculateAllOfTfidf(tfidf *Tfidf) {
	for _, document := range tfidf.data {
		total := document.total()
		for word, value := range document.words {
			newValue := calculateTfidf(tfidf, value.word, value.count, total)
			document.words[word] = newValue
		}
	}
}

func updateLicense(tfidf *Tfidf, license *License) {
	doc := &document{name: license.Name, words: map[string]*value{}}
	for word, count := range license.Frequencies {
		doc.words[word] = &value{word: word, count: count}
	}
	tfidf.data[license.Name] = doc
}

/*
Prepare of tfidf calculating tfidf of each document.
*/
func (tfidf *Tfidf) Prepare(db *Database) error {
	licenses := db.Entries("tfidf")
	for _, license := range licenses {
		updateLicense(tfidf, license)
	}
	calculateAllOfTfidf(tfidf)
	return nil
}

/*
Parse parses given data and create an instance of License by tfidf.
*/
func (tfidf *Tfidf) Parse(reader io.Reader, licenseName string) (*License, error) {
	result, err := readFully(reader)
	if err != nil {
		return nil, err
	}
	return BuildWordFreqLicense(licenseName, result)
}

func extractKeysFromDocument(doc1, doc2 *document) map[string]int {
	keys := map[string]int{}
	for _, val := range doc1.words {
		keys[val.word] = 1
	}
	for _, val := range doc2.words {
		keys[val.word] = 1
	}
	return keys
}

func similarity(doc1, doc2 *document) float64 {
	keys := extractKeysFromDocument(doc1, doc2)
	sum := float64(0)
	for key := range keys {
		sum += (doc1.tfidf(key) * doc2.tfidf(key))
	}
	return sum / (doc1.magnitude() * doc2.magnitude())
}

/*
Compare computes similarity between given two licenses.
*/
func (tfidf *Tfidf) Compare(license1, license2 *License) float64 {
	doc1 := findDocument(tfidf, license1)
	doc2 := findDocument(tfidf, license2)
	return similarity(doc1, doc2)
}

func findDocument(tfidf *Tfidf, license *License) *document {
	doc, ok := tfidf.data[license.Name]
	if ok {
		return doc
	}
	doc = &document{name: license.Name, words: map[string]*value{}}
	total := license.total()
	for word, count := range license.Frequencies {
		doc.words[word] = calculateTfidf(tfidf, word, count, total)
	}
	return doc
}
