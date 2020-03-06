package lioss

import (
	"archive/zip"
	"fmt"
)

/*
ZipProject shows an project formatted in zip file.
*/
type ZipProject struct {
	path       string
	readCloser *zip.ReadCloser
}

/*
Close closes project.
*/
func (zp *ZipProject) Close() error {
	if zp.readCloser != nil {
		return zp.readCloser.Close()
	}
	return nil
}

/*
BasePath returns the path of the project.
*/
func (zp *ZipProject) BasePath() string {
	return zp.path
}

func readFileNames(zp *ZipProject) []string {
	results := []string{}
	for _, file := range zp.readCloser.Reader.File {
		if IsLicenseFile(file.FileHeader.Name) {
			results = append(results, file.FileHeader.Name)
		}
	}
	return results
}

/*
LicenseIDs returns ids containing the project for LicenseFile method.
*/
func (zp *ZipProject) LicenseIDs() []string {
	if zp.readCloser == nil {
		reader, err := zip.OpenReader(zp.BasePath())
		if err != nil {
			return []string{}
		}
		zp.readCloser = reader
	}
	return readFileNames(zp)
}

/*
LicenseFile finds the license file path from project.
*/
func (zp *ZipProject) LicenseFile(licenseID string) (LicenseFile, error) {
	for _, file := range zp.readCloser.Reader.File {
		if file.FileHeader.Name == licenseID {
			reader, err := file.Open()
			return &BasicLicenseFile{id: licenseID, reader: reader}, err
		}
	}
	return nil, fmt.Errorf("%s: not found", licenseID)
}
