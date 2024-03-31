package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/petems/go-sshconfig"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Version is what is returned by the `-v` flag
const Version = "0.1.0"

// gitCommit is the gitcommit its built from
var gitCommit = "development"

var (
	app = kingpin.New("sck", "A command-line tool for modifying ssh config files.")

	debug    = app.Flag("debug", "Enable debug mode.").Bool()
	diff     = app.Flag("diff", "Shows a diff").Default("true").Bool()
	filePath = app.Flag("filepath", "Filepath to the sshconfig").Default("$HOME/.ssh/config").String()

	host           = app.Command("host", "Make changes to a host in the ssh config.")
	hostHostname   = host.Flag("hostname", "The hostname you want to change").Short('h').Required().String()
	hostDryRun     = host.Flag("dry-run", "Do a dry-run and don't override the existing file").Default("false").Bool()
	hostBackup     = host.Flag("backup", "Create a backup of the original config file").Default("true").Bool()
	hostCreate     = host.Flag("create", "Creates new values if not present").Default("true").Bool()
	hostParam      = host.Flag("param", "The param you want to change").Required().String()
	hostParamValue = host.Flag("value", "The value you want to change").Required().String()

	global           = app.Command("global", "Make changes to a global part of the ssh config.")
	globalDryRun     = global.Flag("dry-run", "Do a dry-run and don't override the existing file").Default("false").Bool()
	globalBackup     = global.Flag("backup", "Create a backup of the original config file").Default("true").Bool()
	globalParam      = global.Flag("param", "The param you want to change").Required().String()
	globalParamValue = global.Flag("value", "The value you want to change").Required().String()

	version = app.Command("version", "Show the version.")
)

func main() {
	kingpin.Version("0.0.1")

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case host.FullCommand():
		sshConfigFile := os.ExpandEnv(*filePath)

		file, err := os.Open(sshConfigFile)
		if err != nil {
			kingpin.Fatalf("Error when looking opening ssh config file: %s", err)
		}

		dat, err := os.ReadFile(sshConfigFile)

		if err != nil {
			kingpin.Fatalf("Error when storing ssh config file: %s", err)
		}

		config, err := sshconfig.Parse(file)
		if err != nil {
			kingpin.Fatalf("Error when parsing ssh config file: %s", err)
		}

		file.Close()

		// create new param
		newParam := sshconfig.NewParam(*hostParam, []string{*hostParamValue}, nil)

		// Look for host
		if host := config.FindByHostname(*hostHostname); host != nil {
			host := config.FindByHostname(*hostHostname)
			host.AddParam(newParam)
		} else {
			if !*hostCreate {
				kingpin.Fatalf("Error when looking for host: %s", errors.New("could not find host"))
			} else {
				newHost := sshconfig.NewHost([]string{*hostHostname}, nil)
				config.AddHost(newHost)
				host := config.FindByHostname(*hostHostname)
				host.AddParam(newParam)
			}
		}
		if *hostDryRun {
			fmt.Print("New SSH Config:\n")
			if *diff {
				dmp := diffmatchpatch.New()

				buf := new(bytes.Buffer)

				config.WriteTo(buf)

				diffs := dmp.DiffMain(string(dat), buf.String(), false)
				fmt.Println(dmp.DiffPrettyText(diffs))
			} else {
				config.WriteTo(os.Stdin)
			}

		} else {
			fmt.Printf("Written to %s\n", sshConfigFile)
			if err := config.WriteToFilepath(sshConfigFile); err != nil {
				kingpin.Fatalf("Failed to write to ssh config file: %s", err)
			}
		}
	case global.FullCommand():

	case version.FullCommand():
		fmt.Printf("%s %s\n", Version, gitCommit)
	default:
		os.Exit(0)
	}

}
