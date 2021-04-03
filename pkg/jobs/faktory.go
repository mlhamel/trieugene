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
	client  *faktory.Client
}

func NewFaktoryManager(cfg *config.Config) Manager {
	manager := worker.NewManager()
	manager.ProcessStrictPriorityQueues("high", "medium", "low", "default")

	return &FaktoryManager{cfg: cfg, manager: manager}
}

func (f *FaktoryManager) Register(job Job) error {
	f.manager.Register(job.Kind(), job.Run)
	return nil
}

func (f *FaktoryManager) Perform(job Job, msgs ...*Message) error {
	f.cfg.Logger().Debug().Msgf("Instanciating job %s", job.Kind())

	instance := faktory.NewJob(job.Kind(), msgs)
	client, err := f.faktoryClientInstance()
	if err != nil {
		return err
	}
	return client.Push(instance)
}

func (f *FaktoryManager) Run(ctx context.Context) error {
	f.manager.Run()
	return nil
}

func (f *FaktoryManager) faktoryClientInstance() (*faktory.Client, error) {
	if f.client == nil {
		f.cfg.Logger().Debug().Msg("Opening connection with Faktory")
		client, err := faktory.Open()
		if err != nil {
			f.cfg.Logger().Error().Msgf("Cannot open connection with Faktory: %w", err)
			return nil, err
		}
		f.client = client
	}
	return f.client, nil
}
