package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamada/lioss"
)

type options struct {
	dest     string
	format   string
	helpFlag bool
	args     []string
}

func helpMessage() string {
	return `mkliossdb [OPTIONS] <LICENSE...>
OPTIONS
    -d, --dest <DEST>        specifies the destination file path. Default is 'liossdb.json'
    -h, --help               print this message.
LICENSE
    specifies license files.`
}

func parseOptions(args []string) (*options, error) {
	opts := new(options)
	flags := flag.NewFlagSet("mkliossdb", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message.")
	flags.StringVarP(&opts.format, "format", "f", "json", "specifies the destination file format.")
	flags.StringVarP(&opts.dest, "dest", "d", "liossdb.json", "specifies the destination file path.")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if len(flags.Args()) > 1 {
		opts.args = flags.Args()[1:]
	}
	return opts, nil
}

func (opts *options) destination() string {
	if strings.HasSuffix(opts.dest, "."+opts.format) {
		return opts.dest
	}
	index := strings.LastIndex(opts.dest, ".")
	if index < 0 {
		return opts.dest + "." + opts.format
	}
	return opts.dest[0:index] + "." + opts.format
}

func (opts *options) isHelpFlag() bool {
	return opts.helpFlag || len(opts.args) == 0
}

func readLicense(file string, algo lioss.Comparator) (*lioss.License, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return algo.Parse(reader, filepath.Base(file))
}

func performEach(args []string, comparator string) ([]*lioss.License, error) {
	fmt.Printf(`building database for comparator "%s" ...`, comparator)
	algo, err := lioss.CreateComparator(comparator)
	if err != nil {
		return nil, err
	}
	licenses := []*lioss.License{}
	for _, arg := range args {
		license, err := readLicense(arg, algo)
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, license)
	}
	fmt.Println(`done`)
	return licenses, nil
}

func perform(opts *options) int {
	results := map[string][]*lioss.License{}
	for _, algorithm := range lioss.AvailableAlgorithms {
		licenses, err := performEach(opts.args, algorithm)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		results[algorithm] = licenses
	}
	err := lioss.OutputLiossDB(opts.destination(), results)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return 0
}

func goMain(args []string) int {
	opts, err := parseOptions(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	if opts.isHelpFlag() {
		fmt.Println(helpMessage())
		return 0
	}
	return perform(opts)
}

func main() {
	status := goMain(os.Args)
	os.Exit(status)
}
