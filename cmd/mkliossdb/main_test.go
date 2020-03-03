package main

import (
	"os"
	"testing"

	"github.com/tamada/lioss"
)

func TestRun(t *testing.T) {
	goMain([]string{"mkliossdb", "-d", "../../hoge.json", "../../data/BSD"})
	defer os.Remove("../../hoge.json")

	db, err := lioss.LoadDatabase("../../hoge.json")
	if err != nil {
		t.Errorf("load failed: %s", err.Error())
	}
	if len(db.Data) != 10 {
		t.Errorf("database did not fully outputed")
	}
}

func Example_pritHelp() {
	goMain([]string{})
	// Output:
	// mkliossdb [OPTIONS] <LICENSE...>
	// OPTIONS
	//     -d, --dest <DEST>        specifies the destination file path. Default is 'liossdb.json'
	//     -f, --format <FORMAT>    specifies format. Default is 'json'
	//     -h, --help               print this message.
	// LICENSE
	//     specifies license files.
}

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
