package jobs

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
)

type CsvJob struct {
	cfg      *config.Config
	manager  trieugene.Manager
	storejob trieugene.Job
}

func NewCsvJob(cfg *config.Config, manager trieugene.Manager, storejob trieugene.Job) trieugene.Job {
	return &CsvJob{
		cfg:      cfg,
		manager:  manager,
		storejob: storejob,
	}
}

func (o *CsvJob) Kind() string {
	return "csv-rougecombien"
}

func (o *CsvJob) Perform(ctx context.Context, args ...interface{}) error {
	var results = make(map[int][]scraper.Result)

	httpScraper := scraper.NewHttpScraper(o.cfg)
	parser := scraper.NewParser(o.cfg, httpScraper)

	err := parser.Run(ctx, func(ctx context.Context, result scraper.Result) error {
		if len(results[result.TakenAt.Day()]) <= 0 {
			results[result.TakenAt.Day()] = []scraper.Result{}
		}
		results[result.TakenAt.Day()] = append(results[result.TakenAt.Day()], result)
		return nil
	})

	if err != nil {
		return err
	}

	for _, values := range results {
		var messages []interface{}
		for _, result := range values {
			messages = append(messages, trieugene.Message{
				ID:          result.Sha1(),
				Kind:        o.storejob.Kind(),
				ProcessedAt: result.ScrapedAt.Unix(),
				HappenedAt:  result.TakenAt.Unix(),
				Value:       result.Outflow,
			})
		}
		o.storejob.Run(ctx, messages...)
	}

	return nil
}

func (o *CsvJob) Run(ctx context.Context, args ...interface{}) error {
	if err := o.Perform(ctx, args); err != nil {
		return err
	}
	return nil
}
