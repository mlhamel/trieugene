package jobs

import (
	"context"
	"fmt"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
)

type OverflowJob struct {
	cfg     *config.Config
	store   store.Store
	manager trieugene.Manager
	job     trieugene.Job
}

func NewOverflowjob(cfg *config.Config, store store.Store, manager trieugene.Manager) trieugene.Job {
	return &OverflowJob{
		cfg:     cfg,
		store:   store,
		manager: manager,
		job:     trieugene.NewStoreJob("store-rougecombien", cfg, store),
	}
}

func (o *OverflowJob) Kind() string {
	return "overflow-rougecombien"
}

func (o *OverflowJob) Perform(ctx context.Context, args ...interface{}) error {
	scraper.NewScraper(o.cfg, func(ctx context.Context, result scraper.Result) error {
		return o.manager.Perform(o.job, &trieugene.Message{
			ID:          result.Sha1(),
			Kind:        o.job.Kind(),
			ProcessedAt: result.ScrapedAt.Unix(),
			HappenedAt:  result.TakenAt.Unix(),
			Data:        fmt.Sprintf("%f", result.Outflow),
		})
	}).Run(ctx)
	return nil
}

func (o *OverflowJob) Run(ctx context.Context, args ...interface{}) error {
	if err := o.Perform(ctx, args); err != nil {
		return err
	}
	return nil
}
