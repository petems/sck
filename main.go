package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/petems/go-sshconfig"
	"github.com/urfave/cli/v2"
)

// Version is what is returned by the `-v` flag
const Version = "0.1.0"

// gitCommit is the gitcommit its built from
var gitCommit = "development"

func main() {
	app := &cli.App{
		Name:    "sck",
		Usage:   "A simple cli to configure an existing SSH Config",
		Version: gitCommit + "-" + Version,
		Action: func(c *cli.Context) error {
			err := cmdConfigParam(c)
			return err
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: "$HOME/.ssh/config",
				Usage: "Path to the SSH config file",
			},
			&cli.StringFlag{
				Name:     "host",
				Aliases:  []string{"-h"},
				Required: true,
				Usage:    fmt.Sprintf("The host to lookup"),
			},
			&cli.StringFlag{
				Name:     "parameter",
				Aliases:  []string{"p"},
				Required: true,
				Usage:    fmt.Sprintf("The config entry to change"),
			},
			&cli.StringFlag{
				Name:     "value",
				Aliases:  []string{"V"},
				Required: true,
				Usage:    fmt.Sprintf("The value to assign to te parameter"),
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Value: false,
				Usage: "Path to the SSH config file",
			},
		},
	}
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func cmdConfigParam(ctx *cli.Context) (err error) {
	sshConfigFile := os.ExpandEnv(ctx.String("path"))

	file, err := os.Open(sshConfigFile)
	if err != nil {
		return err
	}

	config, err := sshconfig.Parse(file)
	if err != nil {
		return err
	}

	file.Close()

	// Look for host
	if host := config.FindByHostname(ctx.String("host")); host != nil {
		if param := host.GetParam(ctx.String("config")); param != nil {
			// Add identity key field
			param.Args = []string{ctx.String("value")}
		} else {
			// Add identity key field
			identityKeyParam := sshconfig.NewParam(sshconfig.IdentityFileKeyword, []string{ctx.String("value")}, nil)
			host.Params = append(host.Params, identityKeyParam)
		}
	} else {
		// for now, if it doesn't exist, delete it
		return errors.New("COULD NOT FIND HOST")
	}

	if ctx.Bool("dry-run") {
		fmt.Print("New SSH Config:\n")
		config.WriteTo(os.Stdin)
	} else {
		fmt.Printf("Written to %s", ctx.String("path"))
		if err := config.WriteToFilepath(ctx.String("path")); err != nil {
			return err
		}
	}

	return nil

}
