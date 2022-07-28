package jobs

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	base "github.com/mlhamel/trieugene/pkg/scraper"
)

type ParquetJobKwargs struct {
	Cfg      *config.Config
	Manager  trieugene.Manager
	StoreJob trieugene.Job
	Parser   base.Parser
	Scraper  base.Scraper
}

type ParquetJob struct {
	kwargs *ParquetJobKwargs
}

func NewParquetJob(kwargs *ParquetJobKwargs) trieugene.Job {
	return &ParquetJob{kwargs: kwargs}
}

func (p *ParquetJob) Kind() string {
	return "parquet-rougecombien"
}

func (p *ParquetJob) Perform(ctx context.Context, args ...interface{}) error {
	var results = make(map[int][]base.Result)

	err := p.kwargs.Parser.Run(ctx, func(ctx context.Context, result base.Result) error {
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
				Kind:        p.kwargs.StoreJob.Kind(),
				ProcessedAt: result.ScrapedAt.Unix(),
				HappenedAt:  result.TakenAt.Unix(),
				Value:       result.Outflow,
			})
		}
		p.kwargs.StoreJob.Run(ctx, messages...)
	}

	return nil
}

func (p *ParquetJob) Run(ctx context.Context, args ...interface{}) error {
	if err := p.Perform(ctx, args); err != nil {
		return err
	}
	return nil
}
