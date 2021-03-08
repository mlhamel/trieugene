package apps

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/mlhamel/trieugene/pkg/config"
	rougecombien "github.com/mlhamel/trieugene/services/rougecombien/pkg/apps"
)

type Scheduler struct {
	cfg       *config.Config
	scheduler *gocron.Scheduler
}

func NewScheduler(cfg *config.Config) *Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)
	return &Scheduler{
		cfg:       cfg,
		scheduler: scheduler,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	ctx, cancelFun := context.WithCancel(ctx)
	defer cancelFun()
	_, err := s.scheduler.Every(1).Hour().Do(s.rougecombien(ctx))

	if err != nil {
		return err
	}

	s.loop(ctx)

	return nil
}

func (s *Scheduler) RunDevelopment(ctx context.Context) error {
	ctx, cancelFun := context.WithCancel(ctx)
	defer cancelFun()

	_, err := s.scheduler.Every(1).Minute().Do(s.rougecombienDevelopment(ctx))

	if err != nil {
		return err
	}

	s.loop(ctx)

	return nil
}

func (s *Scheduler) rougecombien(ctx context.Context) func() {
	return func() {
		app := rougecombien.NewRougecombien(s.cfg)
		err := app.Run(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Scheduler) rougecombienDevelopment(ctx context.Context) func() {
	return func() {
		app := rougecombien.NewRougecombien(s.cfg)
		err := app.RunDevelopment(ctx)
		if err != nil {
			panic(err)
		}
	}
}

func (s *Scheduler) loop(ctx context.Context) {
	s.scheduler.StartAsync()

	for {
		select {
		case <-ctx.Done():
			s.scheduler.Stop()
		}
	}
}
