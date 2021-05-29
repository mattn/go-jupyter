package main

import (
	"os"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdCat(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:  c.String("token"),
		Origin: c.String("origin"),
	})

	return client.Cat(c.Args().First(), os.Stdout)
}
