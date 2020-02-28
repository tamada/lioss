package lioss

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

/*
LicenseFile shows the path of license in the project.
*/
type LicenseFile interface {
	ID() string
	Read(p []byte) (int, error)
	Close() error
}

/*
Project shows project containing some licenses.
*/
type Project interface {
	BasePath() string
	LicenseIDs() []string
	LicenseFile(licenseID string) (LicenseFile, error)
}

/*
BasicLicenseFile is an instance of LicenseFile.
*/
type BasicLicenseFile struct {
	id     string
	reader io.ReadCloser
}

/*
BasicProject is an instance of Project.
*/
type BasicProject struct {
	baseDir      string
	licensePaths []string
}

/*
BasePath returns the path of the project.
*/
func (project *BasicProject) BasePath() string {
	return project.baseDir
}

/*
LicenseIDs returns ids containing the project for LicenseFile method.
*/
func (project *BasicProject) LicenseIDs() []string {
	return project.licensePaths
}

/*
LicenseFile finds the license file path from project.
*/
func (project *BasicProject) LicenseFile(licenseID string) (LicenseFile, error) {
	path := filepath.Join(project.BasePath(), licenseID)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &BasicLicenseFile{id: licenseID, reader: file}, nil
}

/*
ID returns id of blf.
*/
func (blf *BasicLicenseFile) ID() string {
	return blf.id
}

/*
Read reads data from license file of blf.
*/
func (blf *BasicLicenseFile) Read(p []byte) (int, error) {
	return blf.reader.Read(p)
}

/*
Close closes the file.
*/
func (blf *BasicLicenseFile) Close() error {
	return blf.reader.Close()
}

/*
NewBasicProject construct an instance of project from file system.
*/
func NewBasicProject(baseDir string) *BasicProject {
	project := &BasicProject{baseDir: baseDir, licensePaths: []string{}}
	findLicenseFile(project)
	return project
}

func isContainOtherWord(fileName string) bool {
	ext := filepath.Ext(fileName)
	return len(fileName) > len(ext)+len("license")
}

func isLicenseFile(path string) bool {
	fileName := strings.ToLower(filepath.Base(path))
	return strings.HasPrefix(fileName, "license") && !isContainOtherWord(fileName)
}

func removeBasePath(basePath, path string) string {
	newPath := path
	if strings.HasPrefix(path, basePath) {
		newPath = strings.Replace(path, basePath, "", -1)
	}
	if strings.HasPrefix(newPath, "/") {
		newPath = newPath[1:]
	}
	return newPath
}

func findLicenseFile(project *BasicProject) {
	stats, err := os.Stat(project.BasePath())
	if err != nil {
		return
	}
	if stats.IsDir() {
		findLicenseFileInDir(project)
	} else {
		project.licensePaths = append(project.licensePaths, "")
	}
}

func findLicenseFileInDir(project *BasicProject) {
	filepath.Walk(project.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isLicenseFile(path) {
			path = removeBasePath(project.baseDir, path)
			project.licensePaths = append(project.licensePaths, path)
		}
		return nil
	})
	sort.Slice(project.licensePaths, func(i, j int) bool {
		return len(project.licensePaths[i]) < len(project.licensePaths[j])
	})
}
