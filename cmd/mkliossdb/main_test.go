package main

import (
	"os"
	"testing"

	"github.com/tamada/lioss"
)

func TestPerformEach(t *testing.T) {
	testdata := []struct {
		args       []string
		comparator string
	}{
		{[]string{}, "unknown-comparator-algorithm"},
		{[]string{"not/exist/file"}, "1gram"},
	}
	for _, td := range testdata {
		_, err := performEach(td.args, td.comparator)
		if err == nil {
			t.Errorf("performEach(%v, %s) should fail", td.args, td.comparator)
		}
	}
}

func TestParseOptionFail(t *testing.T) {
	_, err := parseOptions([]string{"mkliossdb", "--unknown"})
	if err == nil {
		t.Errorf("parseOptions should be fail, because it specifies unknown option")
	}
}

func TestOutputError(t *testing.T) {
	err := lioss.OutputLiossDB("not/exist/dir/hoge.json", map[string][]*lioss.License{})
	if err == nil {
		t.Errorf("dabase write should fail, because not exist dir")
	}
}

func TestRun(t *testing.T) {
	goMain([]string{"mkliossdb", "-d", "../../hoge.json", "../../data/misc/BSD"})
	defer os.Remove("../../hoge.json")

	db, err := lioss.LoadDatabase("../../hoge.json")
	if err != nil {
		t.Errorf("load failed: %s", err.Error())
	}
	if len(db.Data) != 11 {
		t.Errorf("database did not fully outputed")
	}
}

func TestIsHelpFlag(t *testing.T) {
	testdata := []struct {
		args         []string
		wontHelpFlag bool
	}{
		{[]string{"mkliossdb", "-h"}, true},
		{[]string{"mkliossdb"}, true},
		{[]string{"mkliossdb", "../../data"}, false},
	}
	for _, td := range testdata {
		opts, err := parseOptions(td.args)
		if err != nil {
			t.Errorf("parseOptions(%v) parse error: %s", td.args, err.Error())
		}
		if opts.isHelpFlag() != td.wontHelpFlag {
			t.Errorf("opts.isHelpFlag() of parseOptions(%v) did not match, wont %v", td.args, td.wontHelpFlag)
		}
	}
}

func Example_pritHelp() {
	goMain([]string{})
	// Output:
	// mkliossdb [OPTIONS] <LICENSE...>
	// OPTIONS
	//     -d, --dest <DEST>        specifies the destination file path. Default is 'liossdb.json'
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
