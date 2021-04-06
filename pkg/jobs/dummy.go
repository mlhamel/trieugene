package jobs

import (
	"context"

	"github.com/mlhamel/trieugene/pkg/config"
)

type DummyManager struct {
	cfg *config.Config
}

func NewDummyManager(cfg *config.Config) Manager {
	return &DummyManager{}
}

func (f *DummyManager) Register(job Job) error {
	return nil
}

func (d *DummyManager) Perform(job Job, msgs ...*Message) error {
	return nil
}

func (f *DummyManager) Run(ctx context.Context) error {
	return nil
}
