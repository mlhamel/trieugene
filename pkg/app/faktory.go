package app

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/pior/runnable"
)

type Faktory struct {
	cfg     *config.Config
	store   store.Store
	manager jobs.Manager
}

func NewFaktory(cfg *config.Config, store store.Store) runnable.Runnable {
	outflowJob := jobs.NewOutflowJob(cfg, store)

	manager := jobs.NewFaktoryManager(cfg)
	manager.Register(outflowJob)

	return &Faktory{
		cfg:     cfg,
		store:   store,
		manager: manager,
	}
}

func (f *Faktory) Run(ctx context.Context) error {
	return f.manager.Run(ctx)
}