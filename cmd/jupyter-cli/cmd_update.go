package main

import (
	"io"
	"os"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdUpdate(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:  c.String("token"),
		Origin: c.String("origin"),
	})

	var in io.Reader = os.Stdin
	if c.Args().Len() > 1 {
		f, err := os.Open(c.Args().Get(1))
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	}

	return client.Save(c.Args().First(), in)
}
