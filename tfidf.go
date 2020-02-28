package lioss

import "io"

/*
Tfidf is an implementation type of Algorithm.
*/
type Tfidf struct {
}

/*
NewTfidf creates an instance of Tfidf.
*/
func NewTfidf() *Tfidf {
	return new(Tfidf)
}

func (tfidf *Tfidf) String() string {
	return "tfidf"
}

/*
Parse parses given data and create an instance of License by tfidf.
*/
func (tfidf *Tfidf) Parse(reader io.Reader, licenseName string) (*License, error) {
	return nil, nil
}

/*
Compare computes similarity between given two licenses.
*/
func (tfidf *Tfidf) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}
