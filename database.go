package lioss

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type Database struct {
	Data map[string][]*License `json:"algorithms"`
}

func NewDatabase() *Database {
	return new(Database)
}

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

func Load(reader io.Reader) (*Database, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	db := NewDatabase()
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) Entries(algorithmName string) []*License {
	return db.Data[algorithmName]
}

func (db *Database) Entry(algoirthmName, licenseName string) *License {
	entries := db.Entries(algoirthmName)
	for _, entry := range entries {
		if entry.Name == licenseName {
			return entry
		}
	}
	return nil
}

func (db *Database) Put(algorithmName string, license *License) {
	if db.Contains(algorithmName, license.Name) {
		// TODO
	} else {
		items, ok := db.Data[license.Name]
		if ok {
			items = append(items, license)
		} else {
			items = []*License{license}
		}
		db.Data[algorithmName] = items
	}
}

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
