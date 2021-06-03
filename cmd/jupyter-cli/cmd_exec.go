package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdExec(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:   c.String("token"),
		Origin:  c.String("origin"),
		Session: c.String("session"),
	})

	notebook, err := client.Notebook(c.Args().First())
	if err != nil {
		return err
	}

	if c.Bool("update") {
		err = notebook.Update(os.Stdin)
		if err != nil {
			return err
		}
	}

	stdout := colorable.NewColorableStdout()
	stderr := colorable.NewColorableStderr()
	if c.Args().Len() == 1 {
		for _, arg := range notebook.CodeIDs() {
			err = notebook.Exec(arg, stdout, stderr)
			if err != nil {
				return err
			}
		}
	} else {
		for _, arg := range c.Args().Slice()[1:] {
			err = notebook.Exec(arg, stdout, stderr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
