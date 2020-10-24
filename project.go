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
	String() string
}

/*
Project shows project containing some licenses.
*/
type Project interface {
	/* BasePath returns the project path. */
	BasePath() string
	/* Close closes the project, and after call this method, this project is not available. */
	Close() error
	/* LicenseIDs returns the ids for licenses in the project. */
	LicenseIDs() []string
	/* LicenseFile returns the path of the License files in the project. */
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
		return &zipProject{path: path}, nil
	}
	if isLicenseFile(filepath.Base(path)) {
		return &dirProject{baseDir: filepath.Dir(path), licensePaths: []string{filepath.Base(path)}}, nil
	}
	return nil, fmt.Errorf("%s: unknown project format", path)
}

/*
basicLicenseFile is an instance of LicenseFile.
*/
type basicLicenseFile struct {
	id     string
	reader io.ReadCloser
}

/*
ID returns id of blf.
*/
func (blf *basicLicenseFile) ID() string {
	return blf.id
}

/*
Read reads data from license file of blf.
*/
func (blf *basicLicenseFile) Read(p []byte) (int, error) {
	return blf.reader.Read(p)
}

/*
Close closes the file.
*/
func (blf *basicLicenseFile) Close() error {
	return blf.reader.Close()
}

func (blf *basicLicenseFile) String() string {
	return blf.ID()
}

func isContainOtherWord(fileName string) bool {
	ext := filepath.Ext(fileName)
	return len(fileName) > len(ext)+len("license")
}

func isLicenseFile(path string) bool {
	fileName := strings.ToLower(filepath.Base(path))
	return strings.HasPrefix(fileName, "license") && !isContainOtherWord(fileName)
}
