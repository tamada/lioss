package lioss

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
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
	Close() error
	LicenseIDs() []string
	LicenseFile(licenseID string) (LicenseFile, error)
}

/*
NewProject creates an instance of Project.
Acceptable file formats of this function is zip/jar/war file, and directory.
*/
func NewProject(path string) (Project, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return newDirProject(path), nil
	}
	return newFileProject(path)
}

func newFileProject(path string) (Project, error) {
	kind, err := filetype.MatchFile(path)
	if err != nil {
		return nil, err
	}
	if kind.MIME.Value == "application/zip" {
		return &ZipProject{path: path}, nil
	}
	return nil, fmt.Errorf("%s: unknown project format", path)
}

/*
BasicLicenseFile is an instance of LicenseFile.
*/
type BasicLicenseFile struct {
	id     string
	reader io.ReadCloser
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

func isContainOtherWord(fileName string) bool {
	ext := filepath.Ext(fileName)
	return len(fileName) > len(ext)+len("license")
}

/*
IsLicenseFile confirms the given path shows the license file.
*/
func IsLicenseFile(path string) bool {
	fileName := strings.ToLower(filepath.Base(path))
	return strings.HasPrefix(fileName, "license") && !isContainOtherWord(fileName)
}
