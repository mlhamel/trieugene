package apps

import (
	"context"
	"fmt"
	"os"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
	"github.com/pior/runnable"
)

type Rougecombien struct {
	cfg *config.Config
}

func NewRougecombien() *Rougecombien {
	cfg := config.NewConfig()
	return &Rougecombien{cfg: cfg}
}

func (r *Rougecombien) Run(ctx context.Context) error {
	return scraper.NewScraper(r.cfg, r.runOutflowJobUsingGCS).Run(ctx)
}

func (r *Rougecombien) RunDevelopment(ctx context.Context) error {
	run(r.setupDevelopment())
	return scraper.NewScraper(r.cfg, r.runOutflowJobUsingS3).Run(ctx)
}

func (r *Rougecombien) runOutflowJobUsingGCS(ctx context.Context, result scraper.Result) error {
	manager := jobs.NewFaktoryManager(r.cfg)
	store, err := store.NewGoogleCloudStorage(ctx, r.cfg)
	if err != nil {
		return err
	}
	job := jobs.NewOutflowJob(r.cfg, store)

	return manager.Perform(job, jobs.Message{
		ID:          result.Sha1(),
		ProcessedAt: result.ScrapedAt,
		Data:        fmt.Sprintf("%f", result.Outflow),
	})
}

func (r *Rougecombien) runOutflowJobUsingS3(ctx context.Context, result scraper.Result) error {
	r.cfg.Logger().Info().Msg("Running outflow using S3 store.")

	store := store.NewS3(r.cfg)
	manager := jobs.NewFaktoryManager(r.cfg)
	job := jobs.NewOutflowJob(r.cfg, store)

	return manager.Perform(job, jobs.Message{
		ID:          result.Sha1(),
		ProcessedAt: result.ScrapedAt,
		Data:        fmt.Sprintf("%f", result.Outflow),
	})
}

func (r *Rougecombien) setupDevelopment() runnable.Runnable {
	return runnable.Func(func(ctx context.Context) error {
		err := os.Setenv("FAKTORY_URL", r.cfg.FaktoryURL())
		if err != nil {
			return err
		}

		err = os.Setenv("FAKTORY_PROVIDER", "FAKTORY_URL")
		if err != nil {
			return err
		}
		return nil
	})
}

func run(runnables ...runnable.Runnable) {
	runnable.RunGroup(runnables...)
}
