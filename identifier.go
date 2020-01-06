package lioss

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
	algorithm, err := BuildAlgorithm(algorithmName)
	if err != nil {
		return nil, err
	}
	identifier.Algorithm = algorithm
	identifier.BuildMasterLicenses("data")
	return identifier, nil
}

type Algorithm interface {
	Parse(project Project) (*License, error)
	String() string
}

type TfidfAlgorithm struct {
}

type NGramAlgorithm struct {
	ngram int
}

type MasterLicenses struct {
	projects []Project
	licenses map[string]*License
}

func (identifier *Identifier) LicenseOf(project Project) (*License, bool) {
	key := filepath.Base(project.LicensePath())
	value, ok := identifier.master.licenses[key]
	if !ok {
		license, err := identifier.Algorithm.Parse(project)
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
	license, err := identifier.Algorithm.Parse(project)
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

/*
TODO build from some database.
*/
func (identifier *Identifier) BuildMasterLicenses(dbpath string) {
	files, err := ioutil.ReadDir(dbpath)
	projects := []Project{}
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			project := BasicProject{baseDir: dbpath, licenseFile: filepath.Join(dbpath, file.Name())}
			projects = append(projects, &project)
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
	// fmt.Printf("debug: license1: %v, license2: %v\n", license1, license2)
	keys := collectKeys([]string{}, license1.frequencies)
	keys = collectKeys(keys, license2.frequencies)

	numerator := 0
	for _, key := range keys {
		numerator = numerator + license1.Of(key)*license2.Of(key)
	}
	return float64(numerator) / (license1.Magnitude() * license2.Magnitude())
}

func BuildAlgorithm(name string) (Algorithm, error) {
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, "gram") {
		nString := strings.Replace(lowerName, "gram", "", -1)
		value, err := strconv.Atoi(nString)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid algorithm name, %s", name, err.Error())
		}
		return NewNGramAlgorithm(value), nil
	} else if lowerName == "tfidf" {
		return NewTfidfAlgorithm(), nil
	}
	return nil, fmt.Errorf("%s: unknown algorithm", lowerName)
}

func NewTfidfAlgorithm() *TfidfAlgorithm {
	return new(TfidfAlgorithm)
}

func (algo *TfidfAlgorithm) String() string {
	return "tfidf"
}

func (algo *TfidfAlgorithm) Parse(project Project) (*License, error) {
	return nil, nil
}

func NewNGramAlgorithm(n int) *NGramAlgorithm {
	ngram := new(NGramAlgorithm)
	ngram.ngram = n
	return ngram
}

func (algo *NGramAlgorithm) String() string {
	return fmt.Sprintf("%dgram", algo.ngram)
}

func (algo *NGramAlgorithm) Parse(project Project) (*License, error) {
	result, err := readFully(project)
	if err != nil {
		return nil, err
	}
	return algo.buildNGram(result)
}

func (algo *NGramAlgorithm) buildNGram(result string) (*License, error) {
	freq := map[string]int{}
	n := algo.ngram
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
	return NewLicense(algo.String(), freq), nil
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
