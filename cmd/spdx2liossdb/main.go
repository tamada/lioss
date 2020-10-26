package main

import (
	"encoding/json"
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
    -d, --dest <DEST>             specifies the destination.
        --with-deprecated         includes deprecated license.
        --without-deprecated      excludes deprecated license.
        --with-osi-approved       includes OSI approved licenses.
        --without-osi-approved    excludes OSI approved licenses.
    -v, --verbose                 verbose mode.
    -h, --help                    prints this message.
ARGUMENT
    the directory contains SPDX license xml files.
NOTE
    this is the internal command, and will not be distributed to the users.`, prog)
}

type cliOptions struct {
	dest        string
	runtimeOpts *runtimeOptions
	helpFlag    bool
	target      string
}

type withWithout struct {
	with    bool
	without bool
}

type runtimeOptions struct {
	verboseOpt  bool
	osiApproved *withWithout
	deprecated  *withWithout
}

type LicenseData struct {
	meta    *lib.LicenseMeta
	content string
}

func (ro *runtimeOptions) verbose(message string) {
	if ro.verboseOpt {
		fmt.Println(message)
	}
}

func (ro *runtimeOptions) verbosef(format string, v ...interface{}) {
	if ro.verboseOpt {
		fmt.Printf(format, v...)
	}
}

func (ww *withWithout) is() bool {
	return ww.with && !ww.without
}

func (ww *withWithout) String() string {
	if ww.is() {
		return "with"
	}
	return "without"
}

func (ww *withWithout) validate() error {
	if ww.with && ww.without {
		return fmt.Errorf("with and without both options cannot be specified")
	}
	if !ww.with && !ww.without {
		return fmt.Errorf("with and without either option must be specified")
	}
	return nil
}

func isTargetLicense(opts *runtimeOptions, meta *lib.LicenseMeta) bool {
	return isTargetLicenseImpl(opts.deprecated.is(), opts.osiApproved.is(), meta)
}

func isTargetLicenseImpl(deprecated, osiApproved bool, meta *lib.LicenseMeta) bool {
	if deprecated && osiApproved {
		return meta.Deprecated && meta.OsiApproved
	} else if deprecated && !osiApproved {
		return meta.Deprecated && !meta.OsiApproved
	} else if !deprecated && osiApproved {
		return !meta.Deprecated && meta.OsiApproved
	}
	return !meta.Deprecated && !meta.OsiApproved
}

func generateLicense(algo lioss.Algorithm, data *LicenseData, opts *runtimeOptions) (*lioss.License, error) {
	if !isTargetLicense(opts, data.meta) {
		return nil, fmt.Errorf("%s: not target license", data.meta.Names.ShortName)
	}
	return algo.Parse(strings.NewReader(data.content), data.meta.Names.ShortName)
}

func performEachAlgorithm(db *lioss.Database, algo lioss.Algorithm, licenseData []*LicenseData, opts *runtimeOptions) error {
	for _, data := range licenseData {
		license, err := generateLicense(algo, data, opts)
		if err != nil {
			continue
		}
		db.Put(algo.String(), license)
	}
	return nil
}

func readLicenseDatum(target string, info os.FileInfo) (*LicenseData, error) {
	if info.IsDir() {
		return nil, fmt.Errorf("%s: is dir", info.Name())
	}
	meta, data, err := lib.ReadSPDX(filepath.Join(target, info.Name()))
	if err != nil {
		return nil, err
	}
	return &LicenseData{meta: meta, content: data}, nil
}

func readLicenseData(target string, opts *runtimeOptions) ([]*LicenseData, error) {
	infoList, err := ioutil.ReadDir(target)
	if err != nil {
		return nil, err
	}
	results := []*LicenseData{}
	for _, info := range infoList {
		result, err := readLicenseDatum(target, info)
		if err == nil {
			results = append(results, result)
		}
	}
	return results, nil
}

func performEach(db *lioss.Database, algorithmName string, licenseData []*LicenseData, opts *runtimeOptions) error {
	algo, err := lioss.NewAlgorithm(algorithmName)
	if err != nil {
		return err
	}
	opts.verbose(algorithmName)
	return performEachAlgorithm(db, algo, licenseData, opts)
}

func performImpl(db *lioss.Database, licenseData []*LicenseData, opts *runtimeOptions) (int, error) {
	size := 0
	for _, algorithmName := range lioss.AvailableAlgorithms {
		err := performEach(db, algorithmName, licenseData, opts)
		if err != nil {
			return size, err
		}
		size = len(db.Data[algorithmName])
	}
	return size, nil
}

type generator interface {
	Perform(licenseData []*LicenseData) error
}

type jsonGenerator struct {
	dest string
	from string
	opts *runtimeOptions
}

type liossdbGenerator struct {
	dest string
	opts *runtimeOptions
}

func newGenerator(dest, from string, opts *runtimeOptions) generator {
	if strings.HasSuffix(dest, ".json") {
		return &jsonGenerator{dest: dest, from: from, opts: opts}
	}
	return &liossdbGenerator{dest: dest, opts: opts}
}

type jsonData struct {
	Timestamp *lioss.Time        `json:"timestamp"`
	CommitID  string             `json:"git-commit-id"`
	Licenses  []*lib.LicenseMeta `json:"licenses"`
}

func (jg *jsonGenerator) Perform(data []*LicenseData) error {
	id, err := readCommitID(jg.from)
	if err != nil {
		fmt.Printf("readCommitID(\"%s\"): failed, %s\n", jg.from, err.Error())
	}
	results := &jsonData{Timestamp: lioss.Now(), CommitID: id, Licenses: []*lib.LicenseMeta{}}
	for _, datum := range data {
		results.Licenses = append(results.Licenses, datum.meta)
	}
	return jg.writeImpl(results)
}

func (jg *jsonGenerator) writeImpl(results *jsonData) error {
	writer, err := os.OpenFile(jg.dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer writer.Close()
	bytes, err := json.Marshal(results)
	if err != nil {
		return err
	}
	length, err := writer.Write(bytes)
	if err != nil {
		return err
	}
	if length != len(bytes) {
		return fmt.Errorf("cannot write fully data, wont %d bytes, write %d bytes", len(bytes), length)
	}
	return nil
}

func (ldg *liossdbGenerator) Perform(data []*LicenseData) error {
	fmt.Printf("read SPDX licenses %s-osi-approved, and %s-deprecated\n", ldg.opts.osiApproved.String(), ldg.opts.deprecated.String())
	db := lioss.NewDatabase()
	size, err := performImpl(db, data, ldg.opts)
	if err != nil {
		return err
	}
	fmt.Printf("parse %d licenses for %d algorithms, and write database to %s...", size, len(db.Data), ldg.dest)
	err2 := db.WriteTo(ldg.dest)
	fmt.Println(" done")
	return err2
}

func perform(dest, target string, opts *runtimeOptions) error {
	licenseData, err := readLicenseData(target, opts)
	if err != nil {
		return err
	}
	generator := newGenerator(dest, target, opts)
	return generator.Perform(licenseData)
}

func buildFlagSet(args []string) (*flag.FlagSet, *cliOptions) {
	opts := new(cliOptions)
	opts.runtimeOpts = &runtimeOptions{osiApproved: &withWithout{}, deprecated: &withWithout{}}
	flags := flag.NewFlagSet("spdx2liossdb", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(helpMessage(args[0])) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.BoolVar(&opts.runtimeOpts.deprecated.without, "without-deprecated", false, "exclude deprecated licenses")
	flags.BoolVar(&opts.runtimeOpts.osiApproved.without, "without-osi-approved", false, "exclude OSI approved licenses")
	flags.BoolVar(&opts.runtimeOpts.deprecated.with, "with-deprecated", false, "exclude deprecated licenses")
	flags.BoolVar(&opts.runtimeOpts.osiApproved.with, "with-osi-approved", false, "exclude OSI approved licenses")
	flags.BoolVarP(&opts.runtimeOpts.verboseOpt, "verbose", "v", false, "verbose mode")
	flags.StringVarP(&opts.dest, "dest", "d", "default.liossdb", "specifies destination of liossdb")
	return flags, opts
}

func validateWithAndWithout(dest string, opts *runtimeOptions) error {
	if strings.HasSuffix(dest, ".json") {
		return nil
	}
	if err := opts.deprecated.validate(); err != nil {
		return fmt.Errorf("deprecated: %s", err.Error())
	}
	if err := opts.osiApproved.validate(); err != nil {
		return fmt.Errorf("osi-approved: %s", err.Error())
	}
	return nil
}

func validateOptions(opts *cliOptions, flags *flag.FlagSet) (*cliOptions, error) {
	if len(flags.Args()) <= 1 {
		return nil, fmt.Errorf("no arguments specified")
	}
	realArgs := flags.Args()[1:]
	if len(realArgs) > 1 {
		return nil, fmt.Errorf("arguments too much: %v", realArgs)
	}
	if err := validateWithAndWithout(opts.dest, opts.runtimeOpts); err != nil {
		return nil, err
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
