package jobs

import (
	"context"

	faktory "github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"github.com/mlhamel/trieugene/pkg/config"
)

type FaktoryManager struct {
	cfg     *config.Config
	manager *worker.Manager
}

func NewFaktoryManager(cfg *config.Config) Manager {
	manager := worker.NewManager()
	manager.ProcessStrictPriorityQueues("high", "medium", "low")

	return &FaktoryManager{cfg: cfg, manager: manager}
}

func (f *FaktoryManager) Register(job Job) error {
	f.manager.Register(job.Name(), job.Run)
	return nil
}

func (f *FaktoryManager) Perform(job Job, args ...interface{}) error {
	client, err := faktory.Open()
	if err != nil {
		return err
	}

	instance := faktory.NewJob(job.Name(), args...)
	return client.Push(instance)
}

func (f *FaktoryManager) Run(ctx context.Context) error {
	f.manager.Run()
	return nil
}
