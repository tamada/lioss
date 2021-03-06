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
	db := lioss.NewDatabase()
	for _, td := range testdata {
		err := performEach(db, td.args, td.comparator)
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
	db := lioss.NewDatabase()
	err := db.WriteTo("not/exist/dir/hoge.liossdb")
	if err == nil {
		t.Errorf("dabase write should fail, because not exist dir")
	}
}

func TestRun(t *testing.T) {
	goMain([]string{"mkliossdb", "-d", "../../hoge.liossdb", "../../data/misc/BSD"})
	defer os.Remove("../../hoge.liossdb")

	db, err := lioss.ReadDatabase("../../hoge.liossdb")
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
	goMain([]string{"mkliossdb"})
	// Output:
	// mkliossdb [OPTIONS] <LICENSEs...>
	// OPTIONS
	//     -d, --dest <DEST>        specifies the destination file path. Default is 'default.liossdb'
	//     -h, --help               print this message.
	// LICENSE
	//     specifies license files.
}
