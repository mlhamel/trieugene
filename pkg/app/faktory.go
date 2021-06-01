package app

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	rougecombienJobs "github.com/mlhamel/trieugene/services/rougecombien/pkg/jobs"
	rougecombienScraper "github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
	"github.com/pior/runnable"
)

type Faktory struct {
	cfg     *config.Config
	store   store.Store
	manager jobs.Manager
}

func NewFaktory(cfg *config.Config, store store.Store) runnable.Runnable {
	jsonStoreJob := jobs.NewJsonStoreJob("json-store-rougecombien", cfg, store)
	csvStoreJob := jobs.NewCsvStoreJob("csv-store-rougecombien", cfg, store)

	manager := jobs.NewFaktoryManager(cfg)

	httpScraper := rougecombienScraper.NewHttpScraper(cfg)
	parser := rougecombienScraper.NewParser(cfg, httpScraper)

	csvJob := rougecombienJobs.NewCsvJob(&rougecombienJobs.CsvJobKwargs{
		Cfg:      cfg,
		Manager:  manager,
		StoreJob: csvStoreJob,
		Scraper:  httpScraper,
		Parser:   parser,
	})

	manager.Register(jsonStoreJob)
	manager.Register(csvStoreJob)
	manager.Register(rougecombienJobs.NewJsonJob(cfg, manager, jsonStoreJob))
	manager.Register(csvJob)

	return &Faktory{
		cfg:     cfg,
		store:   store,
		manager: manager,
	}
}

func (f *Faktory) Run(ctx context.Context) error {
	return f.manager.Run(ctx)
}
