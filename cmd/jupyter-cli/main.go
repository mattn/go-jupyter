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
		&cli.StringFlag{
			Name:    "session",
			Usage:   "jupyter session ID",
			EnvVars: []string{"JUPYTER_CLI_SESSION"},
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
			Name:  "code",
			Usage: "show code",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "json",
					Usage: "output JSON",
				},
			},
			Action: cmdCode,
		},
		{
			Name:  "doc",
			Usage: "show document",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "json",
					Usage: "output JSON",
				},
			},
			Action: cmdDoc,
		},
		{
			Name:  "exec",
			Usage: "execute content",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "update",
					Usage: "update JSON",
				},
				&cli.BoolFlag{
					Name:  "json",
					Usage: "output JSON",
				},
			},
			Action: cmdExec,
		},
		{
			Name:   "session",
			Usage:  "generate session ID",
			Action: cmdSession,
		},
		{
			Name:   "new-id",
			Usage:  "generate ID",
			Action: cmdNewID,
		},
	}
	app.RunAndExitOnError()
}
