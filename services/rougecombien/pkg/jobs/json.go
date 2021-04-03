package jobs

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
)

type JsonJob struct {
	cfg      *config.Config
	manager  trieugene.Manager
	storejob trieugene.Job
}

func NewJsonJob(cfg *config.Config, manager trieugene.Manager, storejob trieugene.Job) trieugene.Job {
	return &JsonJob{
		cfg:      cfg,
		manager:  manager,
		storejob: storejob,
	}
}

func (o *JsonJob) Kind() string {
	return "json-rougecombien"
}

func (o *JsonJob) Perform(ctx context.Context, args ...interface{}) error {
	httpScraper := scraper.NewHttpScraper(o.cfg)
	parser := scraper.NewParser(o.cfg, httpScraper)

	return parser.Run(ctx, func(ctx context.Context, result scraper.Result) error {
		return o.manager.Perform(o.storejob, &trieugene.Message{
			ID:          result.Sha1(),
			Kind:        o.storejob.Kind(),
			ProcessedAt: result.ScrapedAt.Unix(),
			HappenedAt:  result.TakenAt.Unix(),
			Value:       result.Outflow,
		})
	})
}

func (o *JsonJob) Run(ctx context.Context, args ...interface{}) error {
	if err := o.Perform(ctx, args); err != nil {
		return err
	}
	return nil
}
