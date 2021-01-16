package main

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/app"
	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/pior/runnable"

	"github.com/urfave/cli"
)

func main() {
	cfg := config.NewConfig()

	cliApp := cli.App{
		Name: "trieugene",
		Action: func(*cli.Context) error {
			run(app.NewListener(cfg), app.NewWeb(cfg))
			return nil
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				run(setupDevelopment(cfg), app.NewListener(cfg))
				run(tearDownDevelopment())
				return nil
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(runnables ...runnable.Runnable) {
	runnable.RunGroup(runnables...)
}

func setupDevelopment(cfg *config.Config) runnable.Runnable {
	return runnable.Func(func(ctx context.Context) error {
		err := os.Setenv("STORAGE_EMULATOR_HOST", cfg.GCSURL())
		if err != nil {
			return err
		}

		err = os.Setenv("PUBSUB_EMULATOR_HOST", cfg.PubSubURL())
		if err != nil {
			return err
		}

		err = os.Setenv("PUBSUB_PROJECT_ID", cfg.ProjectID())
		if err != nil {
			return err
		}

		err = os.Setenv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE", "1")
		if err != nil {
			return err
		}

		return nil
	})
}

func tearDownDevelopment() runnable.Runnable {
	return runnable.Func(func(ctx context.Context) error {
		err := os.Unsetenv("STORAGE_EMULATOR_HOST")
		if err != nil {
			return err
		}

		err = os.Unsetenv("PUBSUB_EMULATOR_HOST")
		if err != nil {
			return err
		}

		err = os.Unsetenv("PUBSUB_PROJECT_ID")
		if err != nil {
			return err
		}

		err = os.Unsetenv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE")
		if err != nil {
			return err
		}

		return nil
	})
}
