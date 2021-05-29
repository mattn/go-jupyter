package main

import (
	"os"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdExec(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:  c.String("token"),
		Origin: c.String("origin"),
	})

	notebook, err := client.Notebook(c.Args().First())
	if err != nil {
		return err
	}
	return notebook.Exec(c.Args().Slice()[1:], os.Stdout)
}
