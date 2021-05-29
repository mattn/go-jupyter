package main

import (
	"encoding/json"
	"fmt"
	"os"

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
		if c.Bool("json") {
			err = json.NewEncoder(os.Stdout).Encode(notebook.Doc().Content.Cells)
			if err != nil {
				return err
			}
		} else {
			for _, id := range notebook.CodeIDs() {
				fmt.Println(id)
			}
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
