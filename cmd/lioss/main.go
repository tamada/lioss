package main

import (
        "fmt"
        "os"

        flag "github.com/spf13/pflag"
        "github.com/tamada/uniq2/lib"
)

type options struct {
    helpFlag bool
    dbpath string
    args []string
}

func printHelp(appName string) {
    fmt.Printf(`%s [OPTIONS] <PROJECTS...>
OPTIONS
       --dbpath <DBPATH>    specifying database path.
    -h, --help              print this message.
PROJECTS
    project directories, and/or archive files contains LICENSE file.
`, appName)
}


func buildFlagSet() (*flag.FlagSet, *lib.Options) {
    var opts = new(options)
        var flags = flag.NewFlagSet("lioss", flag.ContinueOnError)
    flags.Usage = func() { printHelp("lioss") }
    flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
    return flags, &opts
}



func goMain(args []string) int {
    return 0
}

func Main() {
    var status = goMain(os.Args)
    os.Exit(status)
}
