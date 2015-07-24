package main

import (
	"github.com/codegangsta/cli"
	"github.com/gofarmer/farm"
	"github.com/gofarmer/utils/config"
)

func main() {
	app := cli.NewApp()
	app.Name = "gofarmer"
	app.Version = config.GetGlobal().Version
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
					Before:    farm.CheckCredentialsConfig,
					Action:    farm.Create,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name: "max-price, m-p",
						},
						cli.StringFlag{
							Name: "max-amount, m-a",
						},
						cli.StringFlag{
							Name: "max-instances, max-i",
						},
						cli.StringFlag{
							Name: "min-instances, min-i",
						},
					},
				},
				{
					Name:      "list",
					ShortName: "l",
					Usage:     "List all farms",
					Action:    farm.List,
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
	app.RunAndExitOnError()
}
