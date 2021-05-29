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
	if c.Args().Len() == 1 {
		for _, arg := range notebook.CodeIDs() {
			err = notebook.Exec(arg, os.Stdout)
			if err != nil {
				return err
			}
		}
	} else {
		for _, arg := range c.Args().Slice()[1:] {
			err = notebook.Exec(arg, os.Stdout)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
