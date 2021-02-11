package main

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/app"
	"github.com/urfave/cli"
)

func main() {
	trieugene := app.NewTrieugene()
	cliApp := cli.App{
		Name: "trieugene",
		Action: func(*cli.Context) error {
			return trieugene.Run(context.Background())
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				return trieugene.RunDevelopment(context.Background())
			},
		},
		{
			Name: "store",
			Action: func(c *cli.Context) error {
				return trieugene.RunStore(context.Background(), c.String("key"), c.String("value"))
			},
			Flags: []cli.Flag{
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
