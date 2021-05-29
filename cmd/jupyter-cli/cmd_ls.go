package main

import (
	"fmt"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdLs(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:  c.String("token"),
		Origin: c.String("origin"),
	})

	items, err := client.List(c.Args().First())
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.Type == "directory" {
			fmt.Println(item.Path + "/")
		} else {
			fmt.Println(item.Path)
		}
	}
	return nil
}
