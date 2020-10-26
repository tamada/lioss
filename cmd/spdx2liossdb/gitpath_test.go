package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindGitRepository(t *testing.T) {
	wd, _ := os.Getwd()
	testdata := []struct {
		givePath string
		wontPath string
	}{
		{"../..", "../../.git"},
		{".", filepath.Clean(filepath.Join(wd, "../..", ".git"))},
		{"../../spdx", "../../.git/modules/spdx"},
		{"../../spdx/src", "../../.git/modules/spdx"},
	}

	for _, td := range testdata {
		gotPath, err := findGitRepository(td.givePath)
		if err != nil {
			t.Fatal(err)
		}
		if gotPath != td.wontPath {
			t.Errorf(`the result of findGitRepository("%s") did not match, wont %s, got %s`, td.givePath, td.wontPath, gotPath)
		}
	}
}
