package app

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/store"

	worker "github.com/contribsys/faktory_worker_go"
	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/jobs"
	"github.com/pior/runnable"
)

type Faktory struct {
	cfg     *config.Config
	manager *worker.Manager
	store   store.Store
}

func NewFaktory(cfg *config.Config, store store.Store) runnable.Runnable {
	outflowJob := jobs.NewOutflowJob(cfg, store)

	manager := worker.NewManager()
	manager.ProcessStrictPriorityQueues("high", "medium", "low")
	manager.Register(outflowJob.Name(), outflowJob.Run)

	return &Faktory{
		cfg:     cfg,
		manager: manager,
		store:   store,
	}
}

func (f *Faktory) Run(ctx context.Context) error {
	f.manager.Run()
	return nil
}
