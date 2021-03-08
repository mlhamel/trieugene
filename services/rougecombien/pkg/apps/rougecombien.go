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
	cfg      *config.Config
	storeJob trieugene.Job
	manager  trieugene.Manager
}

func NewRougecombien(cfg *config.Config) *Rougecombien {
	store := store.NewS3(cfg)

	return &Rougecombien{
		cfg:      cfg,
		storeJob: trieugene.NewStoreJob("store-rougecombien", cfg, store),
		manager:  trieugene.NewFaktoryManager(cfg),
	}
}

func (r *Rougecombien) Run(ctx context.Context) error {
	return r.manager.Perform(jobs.NewOverflowjob(r.cfg, r.manager, r.storeJob), &trieugene.Message{})
}

func (r *Rougecombien) RunDevelopment(ctx context.Context) error {
	run(r.setupDevelopment())

	return r.manager.Perform(jobs.NewOverflowjob(r.cfg, r.manager, r.storeJob), &trieugene.Message{})
}

func (r *Rougecombien) genericRun(ctx context.Context, result scraper.Result) error {
	return r.manager.Perform(jobs.NewOverflowjob(r.cfg, r.manager, r.storeJob), &trieugene.Message{})
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
