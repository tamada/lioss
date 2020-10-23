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

func isValidDBType(name string) bool {
	validItems := []string{"whole", "osi", "deprecated", "base"}
	lower := strings.ToLower(name)
	for _, item := range validItems {
		if item == lower {
			return true
		}
	}
	return false
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isValidDBPath(dbpath string) bool {
	if dbpath != "" {
		return existsFile(dbpath)
	}
	return true
}

func validateOptions(opts *liossOptions, args []string) error {
	if !isValidAlgorithm(opts.algorithm) {
		return fmt.Errorf("%s: unknown algorithm", opts.algorithm)
	}
	if !isValidThreshold(opts.threshold) {
		return fmt.Errorf("%f: threshold must be 0.0 to 1.0", opts.threshold)
	}
	if !isValidArgs(args) {
		return fmt.Errorf("no arguments")
	}
	if !isValidDBPath(opts.dbPath) {
		return fmt.Errorf("%s: file not found", opts.dbPath)
	}
	if !isValidDBType(opts.dbtype) {
		return fmt.Errorf("%s: invalid database type", opts.dbtype)
	}
	return nil
}
