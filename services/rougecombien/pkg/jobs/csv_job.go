package jobs

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	base "github.com/mlhamel/trieugene/pkg/scraper"
)

type CsvJobKwargs struct {
	Cfg      *config.Config
	Manager  trieugene.Manager
	StoreJob trieugene.Job
	Parser   base.Parser
	Scraper  base.Scraper
}

type CsvJob struct {
	kwargs *CsvJobKwargs
}

func NewCsvJob(kwargs *CsvJobKwargs) trieugene.Job {
	return &CsvJob{kwargs: kwargs}
}

func (o *CsvJob) Kind() string {
	return "csv-rougecombien"
}

func (o *CsvJob) Perform(ctx context.Context, args ...interface{}) error {
	var results = make(map[int][]base.Result)

	err := o.kwargs.Parser.Run(ctx, func(ctx context.Context, result base.Result) error {
		if len(results[result.TakenAt.Day()]) <= 0 {
			results[result.TakenAt.Day()] = []base.Result{}
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
				Kind:        o.kwargs.StoreJob.Kind(),
				ProcessedAt: result.ScrapedAt.Unix(),
				HappenedAt:  result.TakenAt.Unix(),
				Value:       result.Outflow,
			})
		}
		o.kwargs.StoreJob.Run(ctx, messages...)
	}

	return nil
}

func (o *CsvJob) Run(ctx context.Context, args ...interface{}) error {
	if err := o.Perform(ctx, args); err != nil {
		return err
	}
	return nil
}
