package main

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/services/rougecombien/pkg/apps"
	"github.com/urfave/cli"
)

func main() {
	rougecombien := apps.NewRougecombien()

	cliApp := cli.App{
		Name: "rougecombien",
		Action: func(*cli.Context) error {
			return rougecombien.Run(context.Background())
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				return rougecombien.RunDevelopment(context.Background())
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
