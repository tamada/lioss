package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamada/lioss"
	"github.com/tamada/lioss/lib"
)

func helpMessage(prog string) string {
	return fmt.Sprintf(`%s [OPTIONS] <ARGUMENT>
OPTIONS
    -d, --dest <DEST>           specifies destination.
        --osi-approved          includes only OSI approved licenses.
        --exclude-deprecated    excludes deprecated license.
    -v, --verbose               verbose mode.
    -h, --help                  print this message.
ARGUMENT
    the directory contains SPDX license xml files.`, prog)
}

type cliOptions struct {
	dest        string
	runtimeOpts *runtimeOptions
	helpFlag    bool
	target      string
}

type runtimeOptions struct {
	verbose            bool
	includeOsiApproved bool
	excludeDeprecated  bool
}

func isIgnoreLicense(opts *runtimeOptions, meta *lib.LicenseMeta) bool {
	if opts.includeOsiApproved && !meta.OsiApproved {
		return true
	}
	if opts.excludeDeprecated && meta.Deprecated {
		return true
	}
	return false
}

func readLicense(algo lioss.Algorithm, path string, opts *runtimeOptions) (*lioss.License, error) {
	meta, licenseData, err := lib.ReadSPDX(path)
	if err != nil {
		return nil, err
	}
	if isIgnoreLicense(opts, meta) {
		return nil, nil
	}
	if opts.verbose {
		fmt.Printf("\t%s\n", meta.String())
	}
	return algo.Parse(strings.NewReader(licenseData), meta.Names.ShortName)
}

func appendLicensesIfNeeded(licenses []*lioss.License, license *lioss.License, err error) []*lioss.License {
	if err == nil && license != nil {
		licenses = append(licenses, license)
	}
	return licenses
}

func readLicenses(algo lioss.Algorithm, target string, opts *runtimeOptions, infoList []os.FileInfo) []*lioss.License {
	licenses := []*lioss.License{}
	for _, info := range infoList {
		if !info.IsDir() {
			license, err := readLicense(algo, filepath.Join(target, info.Name()), opts)
			licenses = appendLicensesIfNeeded(licenses, license, err)
		}
	}
	return licenses
}

func performEachAlgorithm(db *lioss.Database, algo lioss.Algorithm, target string, opts *runtimeOptions) error {
	infoList, err := ioutil.ReadDir(target)
	if err != nil {
		return err
	}
	licenses := readLicenses(algo, target, opts, infoList)
	for _, license := range licenses {
		db.Put(algo.String(), license)
	}
	return nil
}

func performEach(db *lioss.Database, algorithmName, target string, opts *runtimeOptions) error {
	algo, err := lioss.NewAlgorithm(algorithmName)
	if err != nil {
		return err
	}
	if opts.verbose {
		fmt.Println(algorithmName)
	}
	return performEachAlgorithm(db, algo, target, opts)
}

func perform(dest, target string, opts *runtimeOptions) error {
	db := lioss.NewDatabase()
	for _, algorithmName := range lioss.AvailableAlgorithms {
		err := performEach(db, algorithmName, target, opts)
		if err != nil {
			return err
		}
	}
	return db.WriteTo(dest)
}

func buildFlagSet(args []string) (*flag.FlagSet, *cliOptions) {
	opts := new(cliOptions)
	opts.runtimeOpts = new(runtimeOptions)
	flags := flag.NewFlagSet("spdx2liossdb", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage(args[0])) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.BoolVar(&opts.runtimeOpts.excludeDeprecated, "exclude-deprecated", false, "exclude deprecated licenses")
	flags.BoolVar(&opts.runtimeOpts.includeOsiApproved, "osi-approved", false, "includes only OSI approved licenses")
	flags.BoolVarP(&opts.runtimeOpts.verbose, "verbose", "v", false, "verbose mode")
	flags.StringVarP(&opts.dest, "dest", "d", "default.liossdb", "specifies destination of liossdb")
	return flags, opts
}

func validateOptions(opts *cliOptions, flags *flag.FlagSet) (*cliOptions, error) {
	if len(flags.Args()) <= 1 {
		return nil, fmt.Errorf("no arguments specified")
	}
	realArgs := flags.Args()[1:]
	if len(realArgs) > 1 {
		return nil, fmt.Errorf("arguments too much: %v", realArgs)
	}
	opts.target = realArgs[0]
	return opts, nil
}

func parseOptions(args []string) (*cliOptions, error) {
	flags, opts := buildFlagSet(args)
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if opts.helpFlag {
		return opts, nil
	}
	return validateOptions(opts, flags)
}

func printError(err error, status int) int {
	if err != nil {
		fmt.Println(err.Error())
		return status
	}
	return 0
}

func goMain(args []string) int {
	opts, err := parseOptions(args)
	if err != nil {
		return printError(err, 1)
	}
	if opts.helpFlag {
		return printError(fmt.Errorf(helpMessage(args[0])), 0)
	}
	if err := perform(opts.dest, opts.target, opts.runtimeOpts); err != nil {
		return printError(err, 2)
	}
	return 0
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
