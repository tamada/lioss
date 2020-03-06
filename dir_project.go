package lioss

import (
	"os"
	"path/filepath"
	"sort"
)

/*
DirProject is an instance of Project.
*/
type DirProject struct {
	baseDir      string
	licensePaths []string
}

func newDirProject(baseDir string) *DirProject {
	project := &DirProject{baseDir: baseDir, licensePaths: []string{}}
	findLicenseFile(project)
	return project
}

/*
Close closes project.
*/
func (project *DirProject) Close() error {
	return nil
}

/*
BasePath returns the path of the project.
*/
func (project *DirProject) BasePath() string {
	return project.baseDir
}

/*
LicenseIDs returns ids containing the project for LicenseFile method.
*/
func (project *DirProject) LicenseIDs() []string {
	return project.licensePaths
}

/*
LicenseFile finds the license file path from project.
*/
func (project *DirProject) LicenseFile(licenseID string) (LicenseFile, error) {
	path := filepath.Join(project.BasePath(), licenseID)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &BasicLicenseFile{id: licenseID, reader: file}, nil
}

func findLicenseFile(project *DirProject) {
	stats, err := os.Stat(project.BasePath())
	if err != nil {
		return
	}
	if stats.IsDir() {
		findLicenseFileInDir(project)
	}
}

func findLicenseFileInDir(project *DirProject) {
	filepath.Walk(project.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && IsLicenseFile(path) {
			path = removeBasePath(project.baseDir, path)
			project.licensePaths = append(project.licensePaths, path)
		}
		return nil
	})
	sort.Slice(project.licensePaths, func(i, j int) bool {
		return len(project.licensePaths[i]) < len(project.licensePaths[j])
	})
}

func min(slice1, slice2 []string) int {
	length1 := len(slice1)
	length2 := len(slice2)
	if length1 < length2 {
		return length1
	}
	return length2
}

func reverse(slice []string) []string {
	results := []string{}
	for i := len(slice) - 1; i >= 0; i-- {
		results = append(results, slice[i])
	}
	return results
}

func filepathToSlice(originalPath string) []string {
	results := []string{}
	path := originalPath
	for path != "." {
		base := filepath.Base(path)
		path = filepath.Dir(path)
		results = append(results, base)
	}
	return reverse(results)
}

/*
IsSamePath tests given two paths are the same.
*/
func IsSamePath(path1, path2 string) bool {
	slice1 := filepathToSlice(path1)
	slice2 := filepathToSlice(path2)
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func removeBasePath(basePath, targetPath string) string {
	bases := filepathToSlice(basePath)
	paths := filepathToSlice(targetPath)
	length := min(bases, paths)
	index := length

	for i := 0; i < length; i++ {
		if bases[i] != paths[i] {
			index = i
			break
		}
	}
	return filepath.Join(paths[index:]...)
}
