package main

import (
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "jupyter-cli"
	app.Usage = "jupyter cli"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "origin",
			Usage:   "origin URL",
			EnvVars: []string{"JUPYTER_CLI_ORIGIN"},
		},
		&cli.StringFlag{
			Name:    "token",
			Usage:   "jupyter API token",
			EnvVars: []string{"JUPYTER_CLI_TOKEN"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "ls",
			Usage:  "list contents",
			Action: cmdLs,
			Flags:  []cli.Flag{},
		},
		{
			Name:   "cat",
			Usage:  "cat content",
			Action: cmdCat,
		},
		{
			Name:   "update",
			Usage:  "update content",
			Action: cmdUpdate,
		},
		{
			Name:   "exec",
			Usage:  "execute content",
			Action: cmdExec,
		},
	}
	app.RunAndExitOnError()
}
