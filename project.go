package lioss

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Project interface {
	Basedir() string
	LicensePath() string
	Open() (io.ReadCloser, error)
	Close()
}

type BasicProject struct {
	baseDir     string
	licenseFile string
	reader      *os.File
}

func NewBasicProject(baseDir string) *BasicProject {
	project := &BasicProject{baseDir: baseDir}
	findLicenseFile(project)
	return project
}

func (project *BasicProject) Basedir() string {
	return project.baseDir
}

func (project *BasicProject) LicensePath() string {
	return project.licenseFile
}

func (project *BasicProject) Open() (io.ReadCloser, error) {
	if project.licenseFile == "" {
		return nil, fmt.Errorf("license file not found")
	}
	reader, err := os.Open(project.licenseFile)
	if err != nil {
		project.reader = reader
	}
	return reader, err
}

func (project *BasicProject) Close() {
	if project.reader != nil {
		project.reader.Close()
	}
}

func isContainOtherWord(fileName string) bool {
	ext := filepath.Ext(fileName)
	return len(fileName) > len(ext)+len("license")
}

func isLicenseFile(path string) bool {
	fileName := strings.ToLower(filepath.Base(path))

	if strings.HasPrefix(fileName, "license") && !isContainOtherWord(fileName) {
		return true
	}
	return false
}

func findLicenseFile(project *BasicProject) {
	licenseFiles := []string{}
	filepath.Walk(project.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isLicenseFile(path) {
			licenseFiles = append(licenseFiles, path)
		}
		return nil
	})
	sort.Slice(licenseFiles, func(i, j int) bool {
		return len(licenseFiles[i]) < len(licenseFiles[j])
	})
	if len(licenseFiles) >= 1 {
		project.licenseFile = licenseFiles[0]
	}
}
