package lioss

import (
	"fmt"
	"sort"
)

/*
Identifier is a type for identifying the license.
*/
type Identifier struct {
	Threshold  float64
	Comparator Algorithm
	Database   *Database
}

/*
Result shows identified license and its probability.
*/
type Result struct {
	/*Name shows the license name.*/
	Name string
	/*Probability represents the probability of the license, by range of 0.0 to 1.0.*/
	Probability float64
}

func (result *Result) String() string {
	return fmt.Sprintf("%s (%f)", result.Name, result.Probability)
}

/*
NewIdentifier creates an instance of Identifier with the given arguments.
The range of threshold must be from 0.0 to 1.0.
*/
func NewIdentifier(algorithmName string, threshold float64, db *Database) (*Identifier, error) {
	identifier := new(Identifier)
	identifier.Threshold = threshold
	algorithm, err := NewAlgorithm(algorithmName)
	if err != nil {
		return nil, err
	}
	identifier.Comparator = algorithm
	identifier.Database = db
	algorithm.Prepare(db)
	return identifier, nil
}

/*Identify identifies the license of the given project. */
func (identifier *Identifier) Identify(project Project) (map[LicenseFile][]*Result, error) {
	ids := project.LicenseIDs()
	resultMap := map[LicenseFile][]*Result{}
	for _, id := range ids {
		file, results, err := identifier.identifyEach(project, id)
		if err != nil {
			return resultMap, err
		}
		resultMap[file] = results
	}
	return resultMap, nil
}

func (identifier *Identifier) identifyEach(project Project, id string) (LicenseFile, []*Result, error) {
	file, err := project.LicenseFile(id)
	if err != nil {
		return file, nil, err
	}
	license, err := identifier.readLicense(file)
	if err != nil {
		return file, nil, err
	}
	results, err := identifier.identify(license)
	return file, results, err
}

/*
ReadLicense reads License from given LicenseFile.
*/
func (identifier *Identifier) readLicense(file LicenseFile) (*License, error) {
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

func (identifier *Identifier) identify(baseLicense *License) ([]*Result, error) {
	licenses := identifier.Database.Entries(identifier.Comparator.String())
	results := []*Result{}
	for _, license := range licenses {
		similarity := identifier.Comparator.Compare(baseLicense, license)
		results = append(results, &Result{Name: license.Name, Probability: similarity})
	}
	return filter(results, identifier.Threshold), nil
}
