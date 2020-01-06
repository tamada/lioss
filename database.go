package lioss

import "time"

type Database struct {
	data map[string][]*Item
}

type Item struct {
	LicenseName   string
	LoadDate      time.Time
	AlgorithmName string
	Data          *License
}

func NewDatabase() *Database {
	return new(Database)
}

func LoadDatabase(path string) (*Database, error) {
	return NewDatabase(), nil
}

func (db *Database) Put(licenseName string, license *License) {
	if db.Contains(licenseName, license.LicenseName) {
		// TODO
	} else {
		items, ok := db.data[licenseName]
		if ok {
			items = append(items, newItem(licenseName, license))
		} else {
			items = []*Item{newItem(licenseName, license)}
		}
		db.data[licenseName] = items
	}
}

func (db *Database) Contains(licenseName string, algorithmName string) bool {
	items, ok := db.data[licenseName]
	if !ok {
		return false
	}
	for _, item := range items {
		if item.AlgorithmName == algorithmName {
			return true
		}
	}
	return false
}

func newItem(licenseName string, license *License) *Item {
	return &Item{LicenseName: licenseName, Data: license, LoadDate: time.Now()}
}
