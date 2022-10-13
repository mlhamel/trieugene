package apps

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/config"
	trieugene "github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/jobs"
	"github.com/mlhamel/trieugene/services/rougecombien/pkg/scraper"
	"github.com/pior/runnable"
)

type Rougecombien struct {
	cfg             *config.Config
	jsonStoreJob    trieugene.Job
	csvStoreJob     trieugene.Job
	parquetStoreJob trieugene.Job
	manager         trieugene.Manager
}

func NewRougecombien(cfg *config.Config) *Rougecombien {
	store := store.NewS3(cfg, &store.S3Params{
		AccessKey:        cfg.S3AccessKey(),
		SecretKey:        cfg.S3SecretKey(),
		URL:              cfg.S3URL(),
		Bucket:           cfg.S3Bucket(),
		Region:           cfg.S3Region(),
		DisableSSL:       true,
		S3ForcePathStyle: true,
	})

	return &Rougecombien{
		cfg:             cfg,
		jsonStoreJob:    trieugene.NewJsonStoreJob("json-store-rougecombien", cfg, store),
		csvStoreJob:     trieugene.NewCsvStoreJob("csv-store-rougecombien", cfg, store),
		parquetStoreJob: trieugene.NewParquetStoreJob("parquet-store-rougecombien", cfg, store),
		manager:         trieugene.NewFaktoryManager(cfg),
	}
}

func NewRougecombienDev(cfg *config.Config) *Rougecombien {
	store := store.NewLocal(cfg)

	return &Rougecombien{
		cfg:             cfg,
		jsonStoreJob:    trieugene.NewJsonStoreJob("json-store-rougecombien", cfg, store),
		csvStoreJob:     trieugene.NewCsvStoreJob("csv-store-rougecombien", cfg, store),
		parquetStoreJob: trieugene.NewParquetStoreJob("parquet-store-rougecombien", cfg, store),
		manager:         trieugene.NewFaktoryManager(cfg),
	}
}

func (r *Rougecombien) Run(ctx context.Context) error {
	httpScraper := scraper.NewHttpScraper(r.cfg)
	parser := scraper.NewParser(r.cfg, httpScraper)

	return r.manager.Perform(jobs.NewParquetJob(&jobs.ParquetJobKwargs{
		Cfg:      r.cfg,
		Manager:  r.manager,
		StoreJob: r.parquetStoreJob,
		Parser:   parser,
		Scraper:  httpScraper,
	}))
}

func (r *Rougecombien) RunDevelopment(ctx context.Context) error {
	run(r.setupDevelopment())

	httpScraper := scraper.NewHttpScraper(r.cfg)
	parser := scraper.NewParser(r.cfg, httpScraper)

	job := jobs.NewParquetJob(&jobs.ParquetJobKwargs{
		Cfg:      r.cfg,
		Manager:  r.manager,
		StoreJob: r.parquetStoreJob,
		Scraper:  httpScraper,
		Parser:   parser,
	})

	return r.manager.Perform(job, &trieugene.Message{})
}

func (r *Rougecombien) RunInline(ctx context.Context) error {
	run(r.setupDevelopment())

	httpScraper := scraper.NewHttpScraper(r.cfg)
	parser := scraper.NewParser(r.cfg, httpScraper)

	job := jobs.NewParquetJob(&jobs.ParquetJobKwargs{
		Cfg:      r.cfg,
		Manager:  r.manager,
		StoreJob: r.parquetStoreJob,
		Scraper:  httpScraper,
		Parser:   parser,
	})

	return job.Run(ctx, &trieugene.Message{})
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
