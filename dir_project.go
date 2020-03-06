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

func removeBasePath(basePath, path string) string {
	newPath := path
	if filepath.HasPrefix(path, basePath) {
		relPath, err := filepath.Rel(basePath, path)
		if err == nil {
			newPath = relPath
		}
	}
	return newPath
}
