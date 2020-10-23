package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/tamada/lioss"
)

/*
VERSION shows the version of the lioss.
*/
const VERSION = "0.9.0"

const defaultDBPath = "data/SPDX-ALL.liossgz"
const dbpathEnvName = "LIOSS_DBPATH"

type liossOptions struct {
	helpFlag  bool
	dbpath    string
	algorithm string
	threshold float64
}

func helpMessage(appName string) string {
	return fmt.Sprintf(`%s version %s
%s [OPTIONS] <PROJECTS...>
OPTIONS
        --dbpath <DBPATH>          specifying database path.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is 5gram.
                                   Available values are: kgram, wordfreq, and tfidf.
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

func printResults(identifier *lioss.Identifier, project lioss.Project) {
	for _, id := range project.LicenseIDs() {
		results, err := identifyLicense(identifier, project, id)
		if err != nil {
			fmt.Printf("%s/%s: %s\n", project.BasePath(), id, err.Error())
			continue
		}
		printResult(project, id, results)
	}
}

func performEach(identifier *lioss.Identifier, arg string, opts *liossOptions) {
	project, err := lioss.NewProject(arg)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	defer project.Close()
	printResults(identifier, project)
	if len(project.LicenseIDs()) == 0 {
		fmt.Printf("%s: license file not found\n", project.BasePath())
	}
}

func databasePath(dbpath string) string {
	if dbpath == defaultDBPath || dbpath == "" {
		if envValue := os.Getenv(dbpathEnvName); envValue != "" {
			return envValue
		}
	}
	return dbpath
}

func perform(args []string, opts *liossOptions) int {
	db, err := lioss.ReadDatabase(opts.dbpath)
	if err != nil {
		return printErrors(err, 1)
	}
	identifier, err := lioss.NewIdentifier(opts.algorithm, opts.threshold, db)
	if err != nil {
		return printErrors(err, 2)
	}
	for _, arg := range args {
		performEach(identifier, arg, opts)
	}
	return 0
}

func buildFlagSet() (*flag.FlagSet, *liossOptions) {
	var opts = new(liossOptions)
	var flags = flag.NewFlagSet("lioss", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage("lioss")) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.StringVarP(&opts.dbpath, "dbpath", "d", defaultDBPath, "specifies database path")
	flags.StringVarP(&opts.algorithm, "algorithm", "a", "5gram", "specifies algorithm")
	flags.Float64VarP(&opts.threshold, "threshold", "t", 0.75, "specifies threshold")
	return flags, opts
}

func parseOptions(args []string, flags *flag.FlagSet, opts *liossOptions) (int, error) {
	if err := flags.Parse(args); err != nil {
		return 1, err
	}
	if opts.isHelpFlag() {
		return 0, fmt.Errorf("%s", helpMessage(args[0]))
	}
	opts.dbpath = databasePath(opts.dbpath)
	if err := validateOptions(opts, flags.Args()[1:]); err != nil {
		return 2, err
	}
	return 0, nil
}

func (opts *liossOptions) isHelpFlag() bool {
	return opts.helpFlag
}

func goMain(args []string) int {
	flags, opts := buildFlagSet()
	status, err := parseOptions(args, flags, opts)
	if err != nil {
		fmt.Println(err.Error())
		return status
	}
	return perform(flags.Args()[1:], opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
