package lioss

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
Database represents the database for the lioss.
*/
type Database struct {
	Timestamp *Time                 `json:"create-at"`
	Data      map[string][]*License `json:"algorithms"`
}

const DatabasePathEnvName = "LIOSS_DBPATH"

type DatabaseType int

const (
	OSI_APPROVED_DATABASE      DatabaseType = 1
	DEPRECATED_DATABASE        DatabaseType = 2
	OSI_DEPRECATED_DATABASE    DatabaseType = 4
	NONE_OSI_APPROVED_DATABASE DatabaseType = 8
	WHOLE_DATABASE             DatabaseType = OSI_APPROVED_DATABASE | DEPRECATED_DATABASE | OSI_DEPRECATED_DATABASE | NONE_OSI_APPROVED_DATABASE
)

func (dt DatabaseType) IsType(dbType DatabaseType) bool {
	if dt <= 0 {
		return false
	}
	if dbType == WHOLE_DATABASE {
		return dbType == dt
	}
	return dt&dbType == dbType
}

func stringImpl(dt DatabaseType) []string {
	array := []string{}
	typeAndNames := []struct {
		t    DatabaseType
		name string
	}{
		{OSI_APPROVED_DATABASE, "OSI_APPROVED_DATABASE"},
		{NONE_OSI_APPROVED_DATABASE, "NONE_OSI_APPROVED_DATABASE"},
		{OSI_DEPRECATED_DATABASE, "OSI_DEPRECATED_DATABASE"},
		{DEPRECATED_DATABASE, "DEPRECATED_DATABASE"},
	}
	for _, tan := range typeAndNames {
		if dt.IsType(tan.t) {
			array = append(array, tan.name)
		}
	}
	return array
}

func (dt DatabaseType) String() string {
	if dt.IsType(WHOLE_DATABASE) {
		return "WHOLE_DATABASE"
	}
	array := stringImpl(dt)
	if len(array) == 0 {
		return "UNKNOWN"
	}
	return strings.Join(array, ",")
}

type dbTypeAndPath struct {
	dbType DatabaseType
	path   string
}

/*LoadDatabase loads the lioss database from system path.
This function search the following directories.

* ENV['LIOSS_DBPATH']
* /usr/local/opt/lioss/data
* /opt/lioss/data
* ./data

If the directory found, this function loads the Base.liossgz, OSIApproved.liossgz, and Deprecated.liossgz as needed.
*/
func LoadDatabase(databaseTypes DatabaseType) (*Database, error) {
	dir, err := availableDatabaseDir()
	if err != nil {
		return nil, err
	}
	dbTypeAndPaths := []dbTypeAndPath{
		{NONE_OSI_APPROVED_DATABASE, "NoneOSIApproved.liossgz"},
		{OSI_DEPRECATED_DATABASE, "OSIDeprecated.liossgz"},
		{OSI_APPROVED_DATABASE, "OSIApproved.liossgz"},
		{DEPRECATED_DATABASE, "Deprecated.liossgz"},
	}
	db := NewDatabase()
	for _, typeAndPath := range dbTypeAndPaths {
		db = loadAndMergeDB(db, dir, databaseTypes, typeAndPath)
	}
	return db, nil
}

func loadAndMergeDB(db *Database, dir string, dbTypes DatabaseType, tp dbTypeAndPath) *Database {
	if dbTypes.IsType(tp.dbType) {
		db2, err := ReadDatabase(filepath.Join(dir, tp.path))
		if err != nil {
			fmt.Printf("%s/%s: %s\n", dir, tp.path, err.Error())
			return db
		}
		db = db.Merge(db2)
	}
	return db
}

func (db *Database) Merge(other *Database) *Database {
	newDB := NewDatabase()
	newDB.Timestamp = db.Timestamp
	newDB.Data = db.Data
	for key, licenses := range other.Data {
		orig := mergeLicense(newDB.Data[key], licenses)
		newDB.Data[key] = orig
	}
	return newDB
}

func mergeLicense(license1, license2 []*License) []*License {
	for _, l := range license2 {
		found := findLicense(l, license1)
		if !found {
			license1 = append(license1, l)
		}
	}
	return license1
}

func findLicense(license *License, array []*License) bool {
	for _, item := range array {
		if item.Name == license.Name {
			return true
		}
	}
	return false
}

func availableDatabaseDir() (string, error) {
	bases := []string{
		os.Getenv(DatabasePathEnvName),
		"/usr/local/opt/lioss/data",
		"/opt/lioss/data",
		"data",
	}
	for _, base := range bases {
		stat, err := os.Stat(base)
		if err == nil && stat.IsDir() {
			return base, nil
		}
	}
	return "", fmt.Errorf("lioss database not found.")
}

/*
NewDatabase create an instance of database for lioss.
*/
func NewDatabase() *Database {
	return &Database{Timestamp: Now(), Data: map[string][]*License{}}
}

func (db *Database) AlgorithmCount() int {
	return len(db.Data)
}

func (db *Database) LicenseCount() int {
	size := 0
	for k, v := range db.Data {
		current := len(v)
		if size != 0 && size != current {
			fmt.Printf("%s: license count not match, size: %d, current: %d\n", k, size, current)
		}
		size = current
	}
	return size
}

/*WriteTo writes data in the the receiver database into the given file.*/
func (db *Database) WriteTo(destFile string) error {
	dest := destination(destFile)
	writer, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()

	newWriter := wrapWriter(writer, destFile)
	err2 := db.Write(newWriter)
	newWriter.Close() // gzip.Writer should call Close.
	writer.Close()
	return err2
}

/*
Write writes database to given writer.
*/
func (db *Database) Write(writer io.Writer) error {
	bytes, err := json.Marshal(db)
	if err != nil {
		return err
	}
	length, err := writer.Write(bytes)
	if err != nil {
		return err
	}
	if length != len(bytes) {
		return fmt.Errorf("cannot write fully data")
	}
	return nil
}

func replaceExtension(fileName, newExt string) string {
	index := strings.LastIndex(fileName, ".")
	if index < 0 {
		return fileName + "." + newExt
	}
	return fileName[0:index] + "." + newExt
}

func destination(dest string) string {
	if strings.HasSuffix(dest, ".liossdb") || strings.HasSuffix(dest, ".liossgz") {
		return dest
	}
	if strings.HasSuffix(dest, ".liossdb.gz") {
		return strings.ReplaceAll(dest, ".liossdb.gz", ".liossgz")
	}
	if strings.HasSuffix(dest, ".gz") {
		return strings.ReplaceAll(dest, ".gz", ".liossgz")
	}
	return replaceExtension(dest, "liossdb")
}

func wrapWriter(writer io.WriteCloser, dest string) io.WriteCloser {
	if strings.HasSuffix(dest, ".gz") || strings.HasSuffix(dest, ".liossgz") {
		return gzip.NewWriter(writer)
	}
	return writer
}

/*
ReadDatabase reads database from given path.
*/
func ReadDatabase(path string) (*Database, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	newReader, err := wrapReader(reader, path)
	if err != nil {
		return nil, err
	}
	return Read(newReader, path)
}

func updateHeader(header gzip.Header, name string) gzip.Header {
	newName := strings.ReplaceAll(name, ".gz", "")
	newName = strings.ReplaceAll(name, ".liossgz", ".liossdb")
	header.Name = newName
	return header
}

func wrapReader(reader io.Reader, from string) (io.Reader, error) {
	if strings.HasSuffix(from, ".liossgz") || strings.HasSuffix(from, ".gz") {
		gz, err := gzip.NewReader(reader)
		gz.Header = updateHeader(gz.Header, from)
		return gz, err
	}
	return reader, nil
}

/*
Read reads database from given reader.
*/
func Read(reader io.Reader, name string) (*Database, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", name, err.Error())
	}
	db := NewDatabase()
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, fmt.Errorf("%s: %s", name, err.Error())
	}
	return db, nil
}

/*
Entries returns a slice of licenses built by given algorithm.
*/
func (db *Database) Entries(algorithmName string) []*License {
	return db.Data[algorithmName]
}

/*
Entry return an instance of license built by given algorithm with given license name.
*/
func (db *Database) Entry(algoirthmName, licenseName string) *License {
	entries := db.Entries(algoirthmName)
	for _, entry := range entries {
		if entry.Name == licenseName {
			return entry
		}
	}
	return nil
}

func updateIfNeeded(items []*License, license *License) (updateFlag bool) {
	for _, item := range items {
		if item.Name == license.Name {
			item.Frequencies = license.Frequencies
			return true
		}
	}
	return false
}

/*
Put registers the given license to the database.
*/
func (db *Database) Put(algorithmName string, license *License) {
	items, ok := db.Data[algorithmName]
	if !ok {
		items = []*License{}
	}
	if !updateIfNeeded(items, license) {
		items = append(items, license)
	}
	db.Data[algorithmName] = items
}

/*
Contains checks existance with algorithm and license name.
*/
func (db *Database) Contains(algorithmName, licenseName string) bool {
	licenses, ok := db.Data[algorithmName]
	if !ok {
		return false
	}
	for _, license := range licenses {
		if license.Name == licenseName {
			return true
		}
	}
	return false
}
