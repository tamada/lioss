package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
	"github.com/tamada/lioss"
)

type mkliossdbOptions struct {
	dest     string
	format   string
	helpFlag bool
	args     []string
}

func helpMessage() string {
	return `mkliossdb [OPTIONS] <LICENSEs...>
OPTIONS
    -d, --dest <DEST>        specifies the destination file path. Default is 'default.liossdb'
    -h, --help               print this message.
LICENSE
    specifies license files.`
}

func buildFlagSet() (*flag.FlagSet, *mkliossdbOptions) {
	opts := new(mkliossdbOptions)
	flags := flag.NewFlagSet("mkliossdb", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage()) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message.")
	flags.StringVarP(&opts.dest, "dest", "d", "default.liossdb", "specifies the destination file path.")
	return flags, opts
}

func parseOptions(args []string) (*mkliossdbOptions, error) {
	flags, opts := buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if len(flags.Args()) > 1 {
		opts.args = flags.Args()[1:]
	}
	return opts, nil
}

func (opts *mkliossdbOptions) isHelpFlag() bool {
	return opts.helpFlag || len(opts.args) == 0
}

func readLicense(file string, algo lioss.Algorithm) (*lioss.License, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return algo.Parse(reader, filepath.Base(file))
}

func performEach(db *lioss.Database, args []string, algorithmName string) error {
	fmt.Printf(`building database for algorithm "%s"...`, algorithmName)
	algorithm, err := lioss.NewAlgorithm(algorithmName)
	if err != nil {
		return err
	}
	for _, arg := range args {
		license, err := readLicense(arg, algorithm)
		if err != nil {
			return err
		}
		db.Put(algorithmName, license)
	}
	fmt.Println(`done`)
	return nil
}

func buildDatabase(opts *mkliossdbOptions) (*lioss.Database, error) {
	db := lioss.NewDatabase()
	for _, algorithm := range lioss.AvailableAlgorithms {
		err := performEach(db, opts.args, algorithm)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
	return db, nil
}

func perform(opts *mkliossdbOptions) int {
	db, err := buildDatabase(opts)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	err = db.WriteTo(opts.dest)
	if err != nil {
		fmt.Println(err.Error())
		return 3
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
