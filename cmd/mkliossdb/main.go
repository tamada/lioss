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
    -f, --format <FORMAT>    specifies format. Default is 'json'
    -h, --help               print this message.
LICENSE
    specifies license files.`
}

func parseOptions(args []string) (*options, error) {
	opts := new(options)
	flags := flag.NewFlagSet("lioss", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message.")
	flags.StringVarP(&opts.dest, "dest", "d", "liossdb.json", "specifies the destination file path.")
	flags.StringVarP(&opts.format, "format", "f", "json", "specifies the format.")
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

var algorithms = []string{"1gram", "2gram", "3gram", "4gram", "5gram", "6gram", "7gram", "8gram", "9gram", "wordfreq", "tfidf"}

func performEach(args []string, comparator string) ([]*lioss.License, error) {
	fmt.Printf(`building database for comparator "%s" ...`, comparator)
	algo, err := lioss.CreateComparator(comparator)
	if err != nil {
		return nil, err
	}
	licenses := []*lioss.License{}
	for _, arg := range args {
		reader, err := os.Open(arg)
		if err != nil {
			return nil, err
		}
		defer reader.Close()
		license, err := algo.Parse(reader, filepath.Base(arg))
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, license)
	}
	fmt.Println(`done`)
	return licenses, nil
}

func output(opts *options, results map[string][]*lioss.License) error {
	db := lioss.NewDatabase()
	db.Data = results
	writer, err := os.OpenFile(opts.destination(), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	return db.Write(writer)
}

func perform(opts *options) int {
	results := map[string][]*lioss.License{}
	for _, algorithm := range algorithms {
		licenses, err := performEach(opts.args, algorithm)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		results[algorithm] = licenses
	}
	err := output(opts, results)
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
