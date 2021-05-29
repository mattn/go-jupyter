package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdDoc(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:  c.String("token"),
		Origin: c.String("origin"),
	})

	notebook, err := client.Notebook(c.Args().First())
	if err != nil {
		return err
	}
	if c.Bool("json") {
		err = json.NewEncoder(os.Stdout).Encode(notebook.Doc())
		if err != nil {
			return err
		}
	} else {
		for _, cell := range notebook.Doc().Content.Cells {
			if cell.CellType == "code" {
				fmt.Println(cell.Source)
				for _, output := range cell.Outputs {
					if s, ok := output["text"]; ok {
						fmt.Println(s)
					}
				}
			}
		}
	}
	return nil
}
