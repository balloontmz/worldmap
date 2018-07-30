package main

import (
	"github.com/lpisces/worldmap/cmds/boot"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "worldmap"
	app.Usage = "country, state, city tree"

	app.Commands = []cli.Command{
		{
			Name:    "boot",
			Aliases: []string{"b"},
			Usage:   "start",
			Action:  boot.Run,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "debug, d",
					Usage: "debug switch",
				},
				cli.StringFlag{
					Name:  "config, c",
					Usage: "load config file",
				},
				cli.StringFlag{
					Name:  "port, p",
					Usage: "srv port",
				},
				cli.StringFlag{
					Name:  "bind",
					Usage: "srv host",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
