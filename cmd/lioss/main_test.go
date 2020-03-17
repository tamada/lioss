package main

import (
	"os"
	"testing"
)

func TestDatabasePath(t *testing.T) {
	testdata := []struct {
		envPath  string
		givePath string
		wontPath string
	}{
		{"", "testdata/test.liossdb", "testdata/test.liossdb"},
		{"envpath", "data/OSIApproved.liossgz", "envpath"},
		{"envpath", "", "envpath"},
		{"envpath", "argspath", "argspath"},
		{"", "", ""},
	}

	for _, td := range testdata {
		if td.envPath != "" {
			os.Setenv(dbpathEnvName, td.envPath)
		}
		gotPath := databasePath(td.givePath)
		if gotPath != td.wontPath {
			t.Errorf("databasePath(%s) did not match, wont %s, got %s", td.givePath, td.wontPath, gotPath)
		}
		os.Unsetenv(dbpathEnvName)
	}
}

func TestInvalidOptions(t *testing.T) {
	testdata := []struct {
		args       []string
		errorFlag  bool
		wontStatus int
		message    string
	}{
		{[]string{"lioss"}, true, 2, "no arguments"},
		{[]string{"lioss", "-a", "unknown"}, true, 2, "unknown: unknown algorithm"},
		{[]string{"lioss", "-t", "2.0"}, true, 2, "2.000000: threshold must be 0.0 to 1.0"},
		{[]string{"lioss", "--dbpath", "no/such/file", "../../LICENSE"}, true, 2, "no/such/file: file not found"},
	}

	for _, td := range testdata {
		flags, opts := buildFlagSet()
		gotStatus, err := parseOptions(td.args, flags, opts)
		if (err != nil) != td.errorFlag {
			t.Errorf("result of parseOptions(%v) did not match, wont error: %v", td.args, td.errorFlag)
		}
		if gotStatus != td.wontStatus {
			t.Errorf("status code of parseOptions(%v) did not match, wont %d, got %d", td.args, td.wontStatus, gotStatus)
		}
		if err != nil && err.Error() != td.message {
			t.Errorf("error message of parseOptions(%v) did not match, wont %s, got %s", td.args, td.message, err.Error())
		}
	}
}

func TestContains(t *testing.T) {
	testdata := []struct {
		item     string
		set      []string
		wontFlag bool
	}{
		{"a", []string{"a", "b", "c"}, true},
		{"b", []string{"a", "b", "c"}, true},
		{"c", []string{"a", "b", "c"}, true},
		{"d", []string{"a", "b", "c"}, false},
		{"abc", []string{"a", "b", "c"}, false},
	}
	for _, td := range testdata {
		gotFlag := contains(td.item, td.set)
		if gotFlag != td.wontFlag {
			t.Errorf("contains(%s, %v), wont %v, got %v", td.item, td.set, td.wontFlag, gotFlag)
		}
	}
}

func Example_invalidDBPath() {
	goMain([]string{"lioss", "--dbpath", "../../testdata/invalid.liossdb", "../../LICENSE"})
	// Output:
	// ../../testdata/invalid.liossdb: unexpected end of JSON input
}

func Example_invalidCLIOptions() {
	goMain([]string{"lioss", "--unknown"})
	// Output:
	// unknown flag: --unknown

}

func Example_lioss() {
	goMain([]string{"lioss", "--dbpath", "../../testdata/test.liossdb", "--algorithm", "6gram", "../../testdata/project3.jar", "../../testdata/project4", "main.go"})
	// Output:
	// ../../testdata/project3.jar/project3/license
	// 	Apache-License-2.0 (1.0000)
	// ../../testdata/project3.jar/project3/subproject/license
	// 	BSD (1.0000)
	// ../../testdata/project4: license file not found
	// main.go: unknown project format
}

func Example_printHelp() {
	goMain([]string{"lioss", "--help"})
	// Output:
	// lioss version 0.9.0
	// lioss [OPTIONS] <PROJECTS...>
	// OPTIONS
	//         --dbpath <DBPATH>          specifying database path.
	//     -a, --algorithm <ALGORITHM>    specifies algorithm. Default is 5gram.
	//                                    Available values are: kgram, wordfreq, and tfidf.
	//     -t, --threshold <THRESHOLD>    specifies threshold of the similarities of license files.
	//                                    Each algorithm has default value. Default value is 0.75.
	//     -h, --help                     print this message.
	// PROJECTS
	//     project directories, and/or archive files contains LICENSE file.
}
