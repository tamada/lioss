package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamada/lioss"
)

/*
VERSION shows the version of the lioss.
*/
const VERSION = "1.0.0"

type options struct {
	helpFlag  bool
	dbpath    string
	algorithm string
	threshold float64
	args      []string
}

func printHelp(appName string) {
	fmt.Printf(`%s version %s
%s [OPTIONS] <PROJECTS...>
OPTIONS
        --dbpath <DBPATH>          specifying database path.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is 5gram.
                                   Available values are: tfidf, kgram, ...
    -t, --threshold <THRESHOLD>    specifies threshold of the similarities of license files.
                                   Each algorithm has default value. Default value is 0.75.
    -h, --help                     print this message.
PROJECTS
    project directories, and/or archive files contains LICENSE file.
`, appName, VERSION, appName)
}

func printResult(project lioss.Project, id string, results []*lioss.Result) {
	fmt.Printf("%s/%s\n", project.BasePath(), id)
	for _, result := range results {
		fmt.Printf("\t%s (%1.4f)\n", result.Name, result.Probability)
	}
}

func printErrors(err error, status int) int {
	fmt.Println(err.Error())
	return status
}

func identifyLicense(identifier *lioss.Identifier, project lioss.Project, id string) ([]*lioss.Result, error) {
	file, err := project.LicenseFile(id)
	if err != nil {
		return nil, err
	}
	license, err := identifier.ReadLicense(file)
	if err != nil {
		return nil, err
	}
	return identifier.Identify(license)
}

func performEach(identifier *lioss.Identifier, arg string, opts *options) {
	project := lioss.NewBasicProject(arg)
	for _, id := range project.LicenseIDs() {
		results, err := identifyLicense(identifier, project, id)
		if err != nil {
			fmt.Printf("%s/%s: %s", project.BasePath(), id, err.Error())
			continue
		}
		printResult(project, id, results)
	}
	if len(project.LicenseIDs()) == 0 {
		fmt.Printf("%s: license file not found\n", project.BasePath())
	}
}

func perform(opts *options) int {
	db, err := lioss.LoadDatabase(opts.dbpath)
	if err != nil {
		return printErrors(err, 1)
	}
	identifier, err := lioss.NewIdentifier(opts.algorithm, opts.threshold, db)
	if err != nil {
		return printErrors(err, 2)
	}
	for _, arg := range opts.args {
		performEach(identifier, arg, opts)
	}
	return 0
}

func buildFlagSet() (*flag.FlagSet, *options) {
	var opts = new(options)
	var flags = flag.NewFlagSet("lioss", flag.ContinueOnError)
	flags.Usage = func() { printHelp("lioss") }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.StringVarP(&opts.dbpath, "dbpath", "d", "liossdb.json", "specifies database path")
	flags.StringVarP(&opts.algorithm, "algorithm", "a", "5gram", "specifies algorithm")
	flags.Float64VarP(&opts.threshold, "threshold", "t", 0.75, "specifies threshold")
	return flags, opts
}

func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func validateOptions(opts *options) error {
	if opts.algorithm != "tfidf" && !strings.HasSuffix(opts.algorithm, "gram") {
		return fmt.Errorf("%s: unknown algorithm", opts.algorithm)
	}
	if opts.threshold < 0.0 || opts.threshold > 1.0 {
		return fmt.Errorf("%f: threshold must be 0.0 to 1.0", opts.threshold)
	}
	if !existsFile(opts.dbpath) {
		return fmt.Errorf("%s: file not found", opts.dbpath)
	}
	return nil
}

func parseOptions(args []string) (*options, error) {
	flags, opts := buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if err := validateOptions(opts); err != nil {
		return opts, err
	}
	opts.args = flags.Args()[1:]
	return opts, nil
}

func (opts *options) isHelpFlag() bool {
	return opts.helpFlag || len(opts.args) == 0
}

func goMain(args []string) int {
	opts, err := parseOptions(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if opts.isHelpFlag() {
		printHelp(args[0])
		return 0
	}
	return perform(opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
