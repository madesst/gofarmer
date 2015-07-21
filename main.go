package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gofarmer/farm"
	"github.com/gofarmer/utils/config"
	"os"
)

func main() {
	fmt.Printf("%+v\n", config.GetGlobal())

	app := cli.NewApp()
	app.Name = "gofarmer"
	app.Usage = "AWS EC2 farm supervisor command-line app"
	app.Author = "madesst"
	app.Email = "madesst@gmail.com"
	app.Commands = []cli.Command{
		{
			Name:      "farm",
			ShortName: "f",
			Usage:     "Operations with farm(s)",
			Subcommands: []cli.Command{
				{
					Name:      "create",
					ShortName: "c",
					Usage:     "Create new farm",
					Action:    farm.Create,
				},
			},
		},
		{
			Name:      "config",
			ShortName: "c",
			Usage:     "Setup global credentials and other stuff",
			Action:    farm.Create,
		},
	}
	app.Run(os.Args)
}
