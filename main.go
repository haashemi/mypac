package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/haashemi/mypac/internal/config"
)

const (
	// prefixSpace is the space used for indentation in the output
	prefixSpace = "  "
)

var (
	boldGreenP  = color.New(color.FgGreen, color.Bold)
	boldRedP    = color.New(color.FgHiRed, color.Bold)
	boldYellowP = color.New(color.FgYellow, color.Bold)
	blueP       = color.New(color.FgBlue)
)

var commands = map[string]func(args []string) error{
	"packages":  runPackages,
	"lostfiles": runLostFiles,
	"modified":  runModified,
}

func main() {
	// No subcommand (or only flags) means `packages`, as it used to be.
	name, args := "packages", os.Args[1:]
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		name, args = args[0], args[1:]
	}

	run, ok := commands[name]
	if !ok {
		usage()
		if name != "help" {
			os.Exit(2)
		}
		return
	}

	if os.Geteuid() != 0 {
		boldRedP.Fprintln(os.Stderr, "mypac must be run as root, try: sudo mypac", strings.Join(os.Args[1:], " "))
		os.Exit(1)
	}

	if err := run(args); err != nil {
		boldRedP.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: sudo mypac [command] [flags]

Commands:
  packages    Categorize explicitly installed packages (default)
  lostfiles   List files not owned by any package, grouped by reason
  modified    List package files modified since installation (mostly /etc)
  help        Show this help

Run 'mypac <command> -h' for the command's flags.

Config search order (unless -config is given):
  %s
`, strings.Join(config.SearchPaths(), "\n  "))
}

// newFlags returns a flag set for a command, with the shared -config flag.
func newFlags(name string) (*flag.FlagSet, *string) {
	fs := flag.NewFlagSet(name, flag.ExitOnError)
	return fs, fs.String("config", "", "path to the config file (default: first one found)")
}
