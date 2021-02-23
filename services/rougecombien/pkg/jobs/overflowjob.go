package jobs

import (
	"context"
	"fmt"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
)

type OverflowJob struct {
	cfg      *config.Config
	manager  trieugene.Manager
	storejob trieugene.Job
}

func NewOverflowjob(cfg *config.Config, manager trieugene.Manager, storejob trieugene.Job) trieugene.Job {
	return &OverflowJob{
		cfg:      cfg,
		manager:  manager,
		storejob: storejob,
	}
}

func (o *OverflowJob) Kind() string {
	return "overflow-rougecombien"
}

func (o *OverflowJob) Perform(ctx context.Context, args ...interface{}) error {
	scraper.NewScraper(o.cfg, func(ctx context.Context, result scraper.Result) error {
		return o.manager.Perform(o.storejob, &trieugene.Message{
			ID:          result.Sha1(),
			Kind:        o.storejob.Kind(),
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
