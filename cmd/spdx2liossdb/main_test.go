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
	//     -d, --dest <DEST>           specifies destination.
	//         --osi-approved          includes only OSI approved licenses.
	//         --exclude-deprecated    excludes deprecated license.
	//     -v, --verbose               verbose mode.
	//     -h, --help                  print this message.
	// ARGUMENT
	//     the directory contains SPDX license xml files.
}

func TestGeneratedDataSize(t *testing.T) {
	testdata := []struct {
		args     []string
		dest     string
		dataSize int
	}{
		{[]string{"spdx2liossdb", "--osi-approved", "--exclude-deprecated", "../../spdx/src"}, "osi.json", 112},
		{[]string{"spdx2liossdb", "--osi-approved", "../../spdx/src"}, "osi_dep.json", 124},
		{[]string{"spdx2liossdb", "--exclude-deprecated", "../../spdx/src"}, "dep.json", 381},
		{[]string{"spdx2liossdb", "../../spdx/src"}, "all.json", 409},
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
	db, err := lioss.LoadDatabase(dest)
	if err != nil {
		t.Errorf("database load error: %s", err.Error())
	}
	if len(db.Data) != len(lioss.AvailableAlgorithms) {
		t.Errorf("data size did not match, wont %d, got %d", len(lioss.AvailableAlgorithms), len(db.Data))
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
		wontIncludeOsiApproved bool
		wontExcludeDeprecated  bool
		wontVerbose            bool
		wontTarget             string
		wontDest               string
	}{
		{[]string{}, true, false, false, false, false, "", ""},
		{[]string{"--unknown-option"}, true, false, false, false, false, "", ""},
		{[]string{"spdx/src"}, false, false, false, false, false, "spdx/src", "liossdb.json"},
		{[]string{"several", "arguments", "causes", "of", "error"}, true, false, false, false, false, "", ""},
		{[]string{"-h"}, false, true, false, false, false, "", "liossdb.json"},
		{[]string{"-v", "-d", "spdx.json", "spdx/src"}, false, false, false, false, true, "spdx/src", "spdx.json"},
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
		if opts.runtimeOpts.verbose != td.wontVerbose {
			t.Errorf("parseOptions(%v) verbose flag did not match, wont %v", td.args, td.wontVerbose)
		}
		if opts.runtimeOpts.excludeDeprecated != td.wontExcludeDeprecated {
			t.Errorf("parseOptions(%v) excludeDeprecated flag did not match, wont %v", td.args, td.wontExcludeDeprecated)
		}
		if opts.runtimeOpts.includeOsiApproved != td.wontIncludeOsiApproved {
			t.Errorf("parseOptions(%v) includeOsiApproved flag did not match, wont %v", td.args, td.wontIncludeOsiApproved)
		}
		if opts.dest != td.wontDest {
			t.Errorf("parseOptions(%v) dest did not match, wont %s, got %s", td.args, td.wontDest, opts.dest)
		}
		if opts.target != td.wontTarget {
			t.Errorf("parseOptions(%v) target did not match, wont %s, got %s", td.args, td.wontTarget, opts.target)
		}
	}
}
