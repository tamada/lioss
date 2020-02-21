package main

import "testing"

func TestUtility(t *testing.T) {
	testdata := []struct {
		giveDest   string
		giveFormat string
		wontDest   string
	}{
		{giveDest: "target.xml", giveFormat: "json", wontDest: "target.json"},
		{giveDest: "target.json", giveFormat: "json", wontDest: "target.json"},
		{giveDest: "target.json.xml", giveFormat: "xml", wontDest: "target.json.xml"},
		{giveDest: "target.", giveFormat: "xml", wontDest: "target.xml"},
		{giveDest: "target", giveFormat: "xml", wontDest: "target.xml"},
	}
	for _, td := range testdata {
		opts := &options{dest: td.giveDest, format: td.giveFormat}
		dest := opts.destination()
		if dest != td.wontDest {
			t.Errorf("destination did not match, wont %s, got %s", td.wontDest, dest)
		}
	}
}
