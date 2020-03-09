package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func contains(word string, set []string) bool {
	for _, item := range set {
		if word == item {
			return true
		}
	}
	return false
}

func isValidAlgorithm(name string) bool {
	if strings.HasSuffix(name, "gram") {
		_, err := strconv.Atoi(strings.ReplaceAll(name, "gram", ""))
		return err == nil
	}
	return contains(name, []string{"tfidf", "wordfreq"})
}

func isValidThreshold(threshold float64) bool {
	return threshold >= 0.0 && threshold <= 1.0
}

func isValidArgs(args []string) bool {
	return len(args) > 0
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isValidDBPath(dbpath string) bool {
	return existsFile(dbpath)
}

func validateOptions(opts *options) error {
	if !isValidAlgorithm(opts.algorithm) {
		return fmt.Errorf("%s: unknown algorithm", opts.algorithm)
	}
	if !isValidThreshold(opts.threshold) {
		return fmt.Errorf("%f: threshold must be 0.0 to 1.0", opts.threshold)
	}
	if !isValidArgs(opts.args) {
		return fmt.Errorf("no arguments")
	}
	if !isValidDBPath(opts.dbpath) {
		return fmt.Errorf("%s: file not found", opts.dbpath)
	}
	return nil
}
