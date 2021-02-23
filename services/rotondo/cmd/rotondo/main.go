package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/services/rotondo/cmd/rotondo/pkg/apps"
	"github.com/urfave/cli"
)

func main() {
	cfg := config.NewConfig()
	scheduler := apps.NewScheduler(cfg)
	ctx := handleSignal(context.Background())

	cliApp := cli.App{
		Name: "rotondo",
		Action: func(*cli.Context) error {
			return scheduler.Run(ctx)
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				return scheduler.RunDevelopment(ctx)
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func handleSignal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx
}
