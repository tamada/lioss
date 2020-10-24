package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamada/lioss"
)

/*
VERSION shows the version of the lioss.
*/
const VERSION = "1.0.0"

type liossOptions struct {
	helpFlag  bool
	dbtype    string
	dbPath    string
	algorithm string
	threshold float64
}

func helpMessage(appName string) string {
	return fmt.Sprintf(`%s version %s
%s [OPTIONS] <PROJECTS...>
OPTIONS
        --database-path <PATH>     specifies the database path.
                                   If specifying this option, database-type option is ignored.
        --database-type <TYPE>     specifies the database type. Default is osi.
                                   Available values are: non-osi, osi, deprecated, osi-deprecated, and whole.
    -a, --algorithm <ALGORITHM>    specifies algorithm. Default is 5gram.
                                   Available values are: kgram, wordfreq, and tfidf.
    -t, --threshold <THRESHOLD>    specifies threshold of the similarities of license files.
                                   Each algorithm has default value. Default value is 0.75.
    -h, --help                     prints this message.
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

func extractKeys(rm map[lioss.LicenseFile][]*lioss.Result) []lioss.LicenseFile {
	slice := []lioss.LicenseFile{}
	for k, _ := range rm {
		slice = append(slice, k)
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].ID() < slice[j].ID()
	})
	return slice
}

func printResults(identifier *lioss.Identifier, project lioss.Project) {
	resultMap, err := identifier.Identify(project)
	if err != nil {
		fmt.Printf(`%s: %s\n`, project.BasePath(), err.Error())
		return
	}
	keys := extractKeys(resultMap)
	for _, key := range keys {
		printResult(project, key.ID(), resultMap[key])
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

func dbType(opts *liossOptions) lioss.DatabaseType {
	switch strings.ToLower(opts.dbtype) {
	case "whole":
		return lioss.WHOLE_DATABASE
	case "osi":
		return lioss.OSI_APPROVED_DATABASE
	case "deprecated":
		return lioss.DEPRECATED_DATABASE
	case "non-osi":
		return lioss.NONE_OSI_APPROVED_DATABASE
	case "osi-deprecated":
		return lioss.OSI_DEPRECATED_DATABASE
	}
	return -1
}

func loadDatabase(opts *liossOptions) (*lioss.Database, error) {
	if opts.dbPath == "" {
		return lioss.LoadDatabase(dbType(opts))
	}
	return lioss.ReadDatabase(opts.dbPath)
}

func perform(args []string, opts *liossOptions) int {
	db, err := loadDatabase(opts)

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
	flags.StringVarP(&opts.algorithm, "algorithm", "a", "5gram", "specifies algorithm")
	flags.StringVarP(&opts.dbtype, "database-type", "d", "osi", "specifies the database type")
	flags.StringVarP(&opts.dbPath, "database-path", "p", "", "specifies the database path")
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
