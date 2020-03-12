package lib

import (
	"compress/gzip"
	"io"
	"os"
	"strings"

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
	newWriter := wrap(writer, dest)
	return db.Write(newWriter)
}

func wrap(writer io.Writer, dest string) io.Writer {
	if strings.HasSuffix(dest, ".gz") {
		return gzip.NewWriter(writer)
	}
	return writer
}
