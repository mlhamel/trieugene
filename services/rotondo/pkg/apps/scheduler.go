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
	s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Running production every hour")

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
	s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Running development every minute")

	ctx, cancelFun := context.WithCancel(ctx)
	defer cancelFun()

	_, err := s.scheduler.Every(1).Minute().Do(s.rougecombienDevelopment(ctx))

	if err != nil {
		s.cfg.Logger().Error().Str("service", "rotondo").Err(err)
		return err
	}

	s.loop(ctx)

	s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Done development every minute")

	return nil
}

func (s *Scheduler) rougecombien(ctx context.Context) func() {
	return func() {
		s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Running rougecombien production")
		app := rougecombien.NewRougecombien(s.cfg)
		err := app.Run(ctx)
		if err != nil {
			s.cfg.Logger().Error().Str("service", "rotondo").Err(err)
			panic(err)
		}
		s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Done rougecombien production")
	}
}

func (s *Scheduler) rougecombienDevelopment(ctx context.Context) func() {
	return func() {
		s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Running rougecombien development")
		app := rougecombien.NewRougecombien(s.cfg)
		err := app.RunDevelopment(ctx)
		if err != nil {
			s.cfg.Logger().Error().Str("service", "rotondo").Err(err)
			panic(err)
		}
		s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Done rougecombien development")
	}
}

func (s *Scheduler) loop(ctx context.Context) {
	s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Running loop")

	s.scheduler.StartBlocking()

	for {
		select {
		case <-ctx.Done():
			s.cfg.Logger().Debug().Str("service", "rotondo").Msg("Stopping scheduler")
			s.scheduler.Stop()
		}
	}
}
