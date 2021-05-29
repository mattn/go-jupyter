package main

import (
	"fmt"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdCode(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:  c.String("token"),
		Origin: c.String("origin"),
	})

	notebook, err := client.Notebook(c.Args().First())
	if err != nil {
		return err
	}
	if c.Args().Len() == 1 {
		ids, err := notebook.CodeIDs()
		if err != nil {
			return err
		}
		for _, id := range ids {
			fmt.Println(id)
		}
	} else {
		for _, arg := range c.Args().Slice() {
			code, err := notebook.Code(arg)
			if err != nil {
				return err
			}
			fmt.Print(code)
		}
	}
	return nil
}