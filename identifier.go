package lioss

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

type Identifier struct {
	Threshold float64
	Algorithm Algorithm
	master    *MasterLicenses
}

type LiossResult struct {
	Name        string
	Probability float64
}

func NewIdentifier(algorithmName string, threshold float64) (*Identifier, error) {
	identifier := new(Identifier)
	identifier.Threshold = threshold
	algorithm, err := CreateAlgorithm(algorithmName)
	if err != nil {
		return nil, err
	}
	identifier.Algorithm = algorithm
	identifier.BuildMasterLicenses("data")
	return identifier, nil
}

type MasterLicenses struct {
	projects []Project
	licenses map[string]*License
}

func (identifier *Identifier) LicenseOf(project Project) (*License, bool) {
	key := filepath.Base(project.LicensePath())
	value, ok := identifier.master.licenses[key]
	if !ok {
		reader, err := project.Open()
		if err != nil {
			return nil, false
		}
		defer reader.Close()
		license, err := identifier.Algorithm.Parse(reader, key)
		if err != nil {
			fmt.Printf("%s: %s\n", project.LicensePath(), err.Error())
			return nil, false
		}
		identifier.master.licenses[key] = license
		value = license
	}
	return value, true
}

func (identifier *Identifier) Identify(project Project) ([]LiossResult, error) {
	reader, err := project.Open()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	license, err := identifier.Algorithm.Parse(reader, filepath.Base(project.LicensePath()))
	if err != nil {
		return nil, err
	}
	results := []LiossResult{}
	for _, masterProject := range identifier.master.projects {
		masterLicense, ok := identifier.LicenseOf(masterProject)
		if !ok {
			fmt.Printf("%s: not found\n", masterProject.LicensePath())
			continue
		}
		similarity := identifier.Compare(masterLicense, license)
		if similarity > identifier.Threshold {
			results = append(results, LiossResult{Probability: similarity, Name: licenseNameOf(masterProject.LicensePath())})
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Probability > results[j].Probability
	})
	return results, nil
}

func licenseNameOf(licensePath string) string {
	return filepath.Base(licensePath)
}

func (identifier *Identifier) BuildMasterLicenses(dbpath string) {
	files, err := ioutil.ReadDir(dbpath)
	projects := []Project{}
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			project := &BasicProject{baseDir: dbpath, licenseFile: filepath.Join(dbpath, file.Name())}
			projects = append(projects, project)
		}
	}
	identifier.master = &MasterLicenses{projects: projects, licenses: map[string]*License{}}
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func collectKeys(keys []string, freq map[string]int) []string {
	for key, _ := range freq {
		if !contains(keys, key) {
			keys = append(keys, key)
		}
	}
	return keys
}

func (identifier *Identifier) Compare(license1, license2 *License) float64 {
	return license1.Similarity(license2)
}
