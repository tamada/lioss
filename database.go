package lioss

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

/*
Database represents the database for the lioss.
*/
type Database struct {
	Data map[string][]*License `json:"algorithms"`
}

/*
NewDatabase create an instance of database for lioss.
*/
func NewDatabase() *Database {
	return &Database{Data: map[string][]*License{}}
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

/*
OutputLiossDB outputs given dbData to file specified in dest.
*/
func OutputLiossDB(dest string, dbData map[string][]*License) error {
	db := NewDatabase()
	db.Data = dbData
	writer, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	newWriter := wrapWriter(writer, dest)
	return db.Write(newWriter)
}

func wrapWriter(writer io.Writer, dest string) io.Writer {
	if strings.HasSuffix(dest, ".gz") {
		return gzip.NewWriter(writer)
	}
	return writer
}

/*
LoadDatabase reads database from given path.
*/
func LoadDatabase(path string) (*Database, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	newReader, err := wrapReader(reader, path)
	if err != nil {
		return nil, err
	}
	return Load(newReader, path)
}

func wrapReader(reader io.Reader, from string) (io.Reader, error) {
	if strings.HasSuffix(from, ".gz") {
		return gzip.NewReader(reader)
	}
	return reader, nil
}

/*
Load reads database from given reader.
*/
func Load(reader io.Reader, name string) (*Database, error) {
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
