package main

import "testing"

func TestInvalidOptions(t *testing.T) {
	testdata := []struct {
		args      []string
		errorFlag bool
		message   string
	}{
		{[]string{"lioss"}, true, "no arguments"},
		{[]string{"lioss", "-a", "unknown"}, true, "unknown: unknown algorithm"},
		{[]string{"lioss", "-t", "2.0"}, true, "2.000000: threshold must be 0.0 to 1.0"},
		{[]string{"lioss", "--dbpath", "no/such/file", "../../LICENSE"}, true, "no/such/file: file not found"},
	}

	for _, td := range testdata {
		_, err := parseOptions(td.args)
		if (err != nil) != td.errorFlag {
			t.Errorf("result of parseOptions(%v) did not match, wont error: %v", td.args, td.errorFlag)
		}
		if err != nil && err.Error() != td.message {
			t.Errorf("error message of parseOptions(%v) did not match, wont %s, got %s", td.args, td.message, err.Error())
		}
	}
}

func Example_main() {
	goMain([]string{"lioss", "--dbpath", "../../testdata/liossdb.json", "--algorithm", "6gram", "../../testdata/project3.jar", "../../testdata/project4"})
	// Output:
	// ../../testdata/project3.jar/project3/license
	// 	Apache-License-2.0 (1.0000)
	// ../../testdata/project3.jar/project3/subproject/license
	// 	BSD (1.0000)
	// ../../testdata/project4: license file not found
}

func Example_printHelp() {
	goMain([]string{"lioss", "--help"})
	// Output:
	// lioss version 1.0.0-beta
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
