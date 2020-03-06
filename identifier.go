package lioss

import (
	"fmt"
	"sort"
)

/*
Identifier is to identify the license.
*/
type Identifier struct {
	Threshold  float64
	Comparator Comparator
	Database   *Database
}

/*
Result shows identified results.
*/
type Result struct {
	Name        string
	Probability float64
}

/*
NewIdentifier creates an instance of Identifier.
*/
func NewIdentifier(comparatorName string, threshold float64, db *Database) (*Identifier, error) {
	identifier := new(Identifier)
	identifier.Threshold = threshold
	algorithm, err := CreateComparator(comparatorName)
	if err != nil {
		return nil, err
	}
	identifier.Comparator = algorithm
	identifier.Database = db
	algorithm.Prepare(db)
	return identifier, nil
}

/*
ReadLicense reads License from given LicenseFile.
*/
func (identifier *Identifier) ReadLicense(file LicenseFile) (*License, error) {
	license, err := identifier.Comparator.Parse(file, file.ID())
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", file.ID(), err.Error())
	}
	return license, nil
}

func filter(results []*Result, threshold float64) []*Result {
	filteredResults := []*Result{}
	for _, r := range results {
		if r.Probability >= threshold {
			filteredResults = append(filteredResults, r)
		}
	}
	sort.Slice(filteredResults, func(i, j int) bool {
		return filteredResults[i].Probability > filteredResults[j].Probability
	})
	return filteredResults
}

/*
Identify identifies the given license.
*/
func (identifier *Identifier) Identify(baseLicense *License) ([]*Result, error) {
	licenses := identifier.Database.Entries(identifier.Comparator.String())
	results := []*Result{}
	for _, license := range licenses {
		similarity := identifier.Comparator.Compare(baseLicense, license)
		results = append(results, &Result{Name: license.Name, Probability: similarity})
	}
	return filter(results, identifier.Threshold), nil
}
