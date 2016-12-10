package main

import (
	"os"

	"github.com/urfave/cli"

	u "fluorescences/utils"
)

func main() {

	app := cli.NewApp()
	app.Name = "Fluorescences"
	app.Usage = "An art gallery blog"
	app.Version = "RC1"
	app.Copyright = "(c) 2016 Tech Janitor"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "tenant",
			Value: "fluorescences",
			Usage: "the name of the tenant",
		},
		cli.StringFlag{
			Name:  "address",
			Value: "localhost",
			Usage: "address to bind to",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 5000,
			Usage: "port to run on",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start the server",
			Action: func(c *cli.Context) error {
				start(c.GlobalString("tenant"), c.GlobalString("address"), c.GlobalInt("port"))
				return nil
			},
		},
		{
			Name:  "init",
			Usage: "initialize a component for the first time",
			Subcommands: []cli.Command{
				{
					Name:  "data",
					Usage: "initialize the boilerplate data",
					Action: func(c *cli.Context) error {
						return u.InitData(c.GlobalString("tenant"))
					},
				},
				{
					Name:  "user",
					Usage: "initialize a user",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "username",
							Usage: "the name of the user",
						},
					},
					Action: func(c *cli.Context) error {
						return u.InitUser(c.GlobalString("tenant"), c.String("username"))
					},
				},
				{
					Name:  "secret",
					Usage: "initialize the HMAC secret",
					Action: func(c *cli.Context) error {
						return u.InitSecret(c.GlobalString("tenant"))
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
