package main

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/apps"
	"github.com/urfave/cli"
)

func main() {
	cfg := config.NewConfig()

	cliApp := cli.App{
		Name: "rougecombien",
		Action: func(*cli.Context) error {
			rougecombien := apps.NewRougecombien(cfg)
			return rougecombien.Run(context.Background())
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				rougecombien := apps.NewRougecombienDev(cfg)
				return rougecombien.RunDevelopment(context.Background())
			},
		},
		{
			Name: "inline",
			Action: func(c *cli.Context) error {
				rougecombien := apps.NewRougecombienDev(cfg)
				return rougecombien.RunInline(context.Background())
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
