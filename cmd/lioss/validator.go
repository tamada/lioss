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

func isValidAlgorithm(opts *liossOptions) error {
	if strings.HasSuffix(opts.algorithm, "gram") {
		_, err := strconv.Atoi(strings.ReplaceAll(opts.algorithm, "gram", ""))
		return err
	}
	if !contains(opts.algorithm, []string{"tfidf", "wordfreq"}) {
		fmt.Errorf("%s: unknown algorithm", opts.algorithm)
	}
	return nil
}

func isValidThreshold(opts *liossOptions) error {
	if opts.threshold >= 0.0 && opts.threshold <= 1.0 {
		return nil
	}
	return fmt.Errorf("%f: threshold must be 0.0 to 1.0", opts.threshold)
}

func isValidArgs(args []string) error {
	if len(args) > 0 {
		return nil
	}
	return fmt.Errorf("no arguments")
}

func isValidDBType(opts *liossOptions) error {
	validItems := []string{"whole", "osi", "deprecated", "base"}
	lower := strings.ToLower(opts.dbtype)
	for _, item := range validItems {
		if item == lower {
			return nil
		}
	}
	return fmt.Errorf("%s: invalid database type", opts.dbtype)
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isValidDBPath(opts *liossOptions) error {
	if opts.dbPath != "" && !existsFile(opts.dbPath) {
		return fmt.Errorf("%s: file not found", opts.dbPath)
	}
	return nil
}

func validateOptions(opts *liossOptions, args []string) error {
	validators := [](func(opts *liossOptions) error){
		isValidAlgorithm, isValidThreshold, isValidDBPath, isValidDBType,
	}
	for _, validator := range validators {
		if err := validator(opts); err != nil {
			return err
		}
	}
	return isValidArgs(args)
}
