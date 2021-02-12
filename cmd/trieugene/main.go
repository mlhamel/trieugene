package main

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/app"
	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/urfave/cli"
)

func main() {
	cfg := config.NewConfig()
	cliApp := cli.App{
		Name: "trieugene",
		Action: func(*cli.Context) error {
			trieugene := app.NewTrieugene(cfg)
			return trieugene.Run(context.Background())
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				trieugene := app.NewTrieugeneDev(cfg)
				return trieugene.Run(context.Background())
			},
		},
		{
			Name: "store",
			Action: func(c *cli.Context) error {
				trieugene := app.NewTrieugeneStore(cfg, c.String("kind"), c.String("key"), c.String("value"))
				return trieugene.Run(context.Background())
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "kind",
					Usage:    "Kind of the stored object",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "key",
					Usage:    "Key of the stored object",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "value",
					Usage:    "value of the stored object",
					Required: true,
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}
