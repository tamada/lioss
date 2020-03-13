package lib

import (
	"testing"
)

func TestReadSPDX(t *testing.T) {
	testdata := []struct {
		path      string
		shortName string
		fullName  string
	}{
		{"../spdx/src/0BSD.xml", "0BSD", "BSD Zero Clause License"},
	}
	for _, td := range testdata {
		meta, _, err := ReadSPDX(td.path)
		if err != nil {
			t.Errorf("ReadSPDX(%s) failed: %s", td.path, err.Error())
		}
		if meta.Names.ShortName != td.shortName || meta.Names.FullName != td.fullName {
			t.Errorf("resultant data of ReadSPDX(%s) did not match, wont %s (%s), got %s", td.path, td.shortName, td.fullName, meta.String())
		}
	}
}
