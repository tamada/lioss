package lib

import (
	"os"

	"github.com/tamada/lioss"
)

func OutputLiossDB(dest string, dbData map[string][]*lioss.License) error {
	db := lioss.NewDatabase()
	db.Data = dbData
	writer, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	return db.Write(writer)
}
