package app

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	rougecombien "github.com/mlhamel/trieugene/services/rougecombien/pkg/jobs"
	"github.com/pior/runnable"
)

type Faktory struct {
	cfg     *config.Config
	store   store.Store
	manager jobs.Manager
}

func NewFaktory(cfg *config.Config, store store.Store) runnable.Runnable {
	storeJob := jobs.NewStoreJob("store-rougecombien", cfg, store)
	manager := jobs.NewFaktoryManager(cfg)
	manager.Register(storeJob)
	manager.Register(rougecombien.NewOverflowjob(cfg, manager, storeJob))

	return &Faktory{
		cfg:     cfg,
		store:   store,
		manager: manager,
	}
}

func (f *Faktory) Run(ctx context.Context) error {
	return f.manager.Run(ctx)
}
