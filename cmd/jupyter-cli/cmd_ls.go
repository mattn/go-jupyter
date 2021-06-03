package main

import (
	"fmt"
	"sort"

	"github.com/mattn/go-jupyter"
	"github.com/urfave/cli/v2"
)

func cmdLs(c *cli.Context) error {
	client := jupyter.NewClient(jupyter.Config{
		Token:   c.String("token"),
		Origin:  c.String("origin"),
		Session: c.String("session"),
	})

	items, err := client.List(c.Args().First())
	if err != nil {
		return err
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Type == "directory" && items[j].Type != "directory" {
			return true
		}
		if items[i].Type != "directory" && items[j].Type == "directory" {
			return false
		}
		if items[i].Name > items[j].Name {
			return false
		}
		return true
	})

	for _, item := range items {
		if item.Type == "directory" {
			fmt.Println(item.Path + "/")
		} else {
			fmt.Println(item.Path)
		}
	}
	return nil
}
