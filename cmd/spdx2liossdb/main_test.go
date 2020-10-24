package main

import (
	"os"
	"sync"
	"testing"

	"github.com/tamada/lioss"
)

func Example_printHelp() {
	goMain([]string{"spdx2liossdb", "--help"})
	// Output:
	// spdx2liossdb [OPTIONS] <ARGUMENT>
	// OPTIONS
	//     -d, --dest <DEST>             specifies the destination.
	//         --with-deprecated         includes deprecated license.
	//         --without-deprecated      excludes deprecated license.
	//         --with-osi-approved       includes OSI approved licenses.
	//         --without-osi-approved    excludes OSI approved licenses.
	//     -v, --verbose                 verbose mode.
	//     -h, --help                    prints this message.
	// ARGUMENT
	//     the directory contains SPDX license xml files.
}

func TestGeneratedDataSize(t *testing.T) {
	testdata := []struct {
		args     []string
		dest     string
		dataSize int
	}{
		{[]string{"spdx2liossdb", "--without-osi-approved", "--without-deprecated", "../../spdx/src"}, "non-osi.liossdb", 269},
		{[]string{"spdx2liossdb", "--with-osi-approved", "--without-deprecated", "../../spdx/src"}, "osi.liossdb", 112},
		{[]string{"spdx2liossdb", "--without-osi-approved", "--with-deprecated", "../../spdx/src"}, "deprecated.liossdb", 16},
		{[]string{"spdx2liossdb", "--with-osi-approved", "--with-deprecated", "../../spdx/src"}, "osi-deprecated.liossdb", 12},
	}

	wg := new(sync.WaitGroup)
	for _, td := range testdata {
		args := append(td.args, "-d")
		args = append(args, td.dest)
		wg.Add(1)
		go testExec(t, args, td.dest, td.dataSize, wg)
	}
	wg.Wait()
}

func testExec(t *testing.T, args []string, dest string, dataSize int, wg *sync.WaitGroup) {
	goMain(args)
	defer os.Remove(dest)
	defer wg.Done()
	db, err := lioss.ReadDatabase(dest)
	if err != nil {
		t.Errorf("testExec(%v): database load error: %s", args, err.Error())
	}
	for key, value := range db.Data {
		if len(value) != dataSize {
			t.Errorf("data size did not match of %s, wont %d, got %d", key, dataSize, len(value))
		}
	}
}

func TestParseOptions(t *testing.T) {
	testdata := []struct {
		args                   []string
		errorFlag              bool
		wontHelp               bool
		wontExcludeOsiApproved bool
		wontExcludeDeprecated  bool
		wontVerbose            bool
		wontTarget             string
		wontDest               string
	}{
		{[]string{}, true, false, false, false, false, "", ""},
		{[]string{"--unknown-option"}, true, false, false, false, false, "", ""},
		{[]string{"spdx/src", "--without-osi-approved", "--without-deprecated"}, false, false, true, true, false, "spdx/src", "default.liossdb"},
		{[]string{"several", "arguments", "causes", "of", "error"}, true, false, false, false, false, "", ""},
		{[]string{"-h"}, false, true, false, false, false, "", "default.liossdb"},
		{[]string{"--with-deprecated", "--without-deprecated"}, true, false, false, false, false, "", "default.liossdb"},
		{[]string{"-v", "-d", "spdx.liossdb", "--with-osi-approved", "--without-deprecated", "spdx/src"}, false, false, false, true, true, "spdx/src", "spdx.liossdb"},
	}

	for _, td := range testdata {
		args := []string{"spdx2liossdb"}
		args = append(args, td.args...)
		opts, err := parseOptions(args)
		if (err != nil) != td.errorFlag {
			t.Errorf("parseOptions(%v) error did not match, wont %v", td.args, td.errorFlag)
		}
		if err != nil {
			continue
		}
		if opts.helpFlag != td.wontHelp {
			t.Errorf("parseOptions(%v) helpFlag did not match, wont %v", td.args, td.wontHelp)
		}
		if opts.runtimeOpts.verboseOpt != td.wontVerbose {
			t.Errorf("parseOptions(%v) verbose flag did not match, wont %v", td.args, td.wontVerbose)
		}
		if opts.runtimeOpts.deprecated.without != td.wontExcludeDeprecated {
			t.Errorf("parseOptions(%v) withoutDeprecated flag did not match, wont %v", td.args, td.wontExcludeDeprecated)
		}
		if opts.runtimeOpts.osiApproved.without != td.wontExcludeOsiApproved {
			t.Errorf("parseOptions(%v) withoutOsiApproved flag did not match, wont %v", td.args, td.wontExcludeOsiApproved)
		}
		if opts.dest != td.wontDest {
			t.Errorf("parseOptions(%v) dest did not match, wont %s, got %s", td.args, td.wontDest, opts.dest)
		}
		if opts.target != td.wontTarget {
			t.Errorf("parseOptions(%v) target did not match, wont %s, got %s", td.args, td.wontTarget, opts.target)
		}
	}
}
