package main

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
	"github.com/pior/runnable"
	"github.com/urfave/cli"
)

func main() {
	cfg := config.NewConfig()

	cliApp := cli.App{
		Name: "rougecombien",
		Action: func(*cli.Context) error {
			ctx := context.Background()
			manager := jobs.NewFaktoryManager(cfg)
			store, err := store.NewGoogleCloudStorage(ctx, cfg)
			if err != nil {
				return err
			}
			run(scraper.NewScraper(cfg, func(ctx context.Context, result *scraper.Result) error {
				job := jobs.NewOutflowJob(cfg, store)
				return manager.Perform(job, result)
			}))
			return nil
		},
	}

	cliApp.Commands = []cli.Command{
		{
			Name: "dev",
			Action: func(c *cli.Context) error {
				store := store.NewS3(cfg)
				manager := jobs.NewFaktoryManager(cfg)
				run(scraper.NewScraper(cfg, func(ctx context.Context, result *scraper.Result) error {
					job := jobs.NewOutflowJob(cfg, store)
					return manager.Perform(job, result)
				}))
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
