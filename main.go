package main

import (
  "fmt"
  "os"
  "github.com/codegangsta/cli"
)

const (
  defaultRunPath    = "/var/run/"
)

func main() {
  app := cli.NewApp()
  app.Name = "gofarmer"
    app.Usage = "sample command-line app"
    app.Author = "madesst"
    app.Email = "madesst@gmail.com"
    app.Commands = []cli.Command{
        {
            Name:      "read",
            ShortName: "r",
            Usage:     "read something",
            Subcommands: []cli.Command{
                {
                    Name:   "tweets",
                    Usage:  "read Tweets",
                    Action: readTweets,
                },
            },
        },
    }
    app.Run(os.Args)
}

func readTweets(ctx *cli.Context) {
    fmt.Println("Go to https://twitter.com/TheProgville to read my tweets!")
}