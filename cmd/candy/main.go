package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/suzuki-shunsuke/candy/pkg/cmd"
	"github.com/suzuki-shunsuke/candy/pkg/domain"
)

var (
	listUpdatedCommand = cli.Command{
		Name:   "list-updated",
		Usage:  "list updated tasks",
		Action: cmd.Run,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "config, c",
				Usage: "configuration file path",
			},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "candy"
	app.Version = domain.Version
	app.Author = "suzuki-shunsuke https://github.com/suzuki-shunsuke"
	app.Usage = "detect updates tasks"
	app.Commands = []cli.Command{
		listUpdatedCommand,
	}
	app.Run(os.Args)
}
