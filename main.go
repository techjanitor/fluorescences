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

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start the server",
			Action: func(c *cli.Context) error {
				name := c.Args().Get(0)
				address := c.Args().Get(1)
				port := c.Args().Get(2)
				start(name, address, port)
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
						name := c.Args().Get(0)
						return u.InitData(name)
					},
				},
				{
					Name:  "user",
					Usage: "initialize a user",
					Action: func(c *cli.Context) error {
						name := c.Args().Get(0)
						user := c.Args().Get(1)
						if user == "" {
							return cli.NewExitError("username required", 1)
						}
						return u.InitUser(name, user)
					},
				},
				{
					Name:  "secret",
					Usage: "initialize the HMAC secret",
					Action: func(c *cli.Context) error {
						name := c.Args().Get(0)
						return u.InitSecret(name)
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
