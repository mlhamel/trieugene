package app

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/pior/runnable"
)

type Trieugene struct {
	cfg *config.Config
}

func NewTrieugene() *Trieugene {
	cfg := config.NewConfig()
	return &Trieugene{cfg: cfg}
}

func (t *Trieugene) Run(ctx context.Context) error {
	store, err := store.NewGoogleCloudStorage(ctx, t.cfg)
	if err != nil {
		return err
	}
	return NewFaktory(t.cfg, store).Run(ctx)
}

func (t *Trieugene) RunDevelopment(ctx context.Context) error {
	store := store.NewS3(t.cfg)

	run(setupDevelopment(t.cfg))

	err := NewFaktory(t.cfg, store).Run(ctx)

	run(tearDownDevelopment())

	return err
}

func (t *Trieugene) RunStore(ctx context.Context, key string, value string) error {
	store := store.NewS3(t.cfg)

	run(setupDevelopment(t.cfg))

	err := NewStore(t.cfg, store, key, value).Run(ctx)

	run(tearDownDevelopment())

	return err
}

func run(runnables ...runnable.Runnable) {
	runnable.RunGroup(runnables...)
}

func setupDevelopment(cfg *config.Config) runnable.Runnable {
	return runnable.Func(func(ctx context.Context) error {
		err := os.Setenv("STORAGE_EMULATOR_HOST", cfg.GCSURL())
		if err != nil {
			return err
		}

		err = os.Setenv("PUBSUB_EMULATOR_HOST", cfg.PubSubURL())
		if err != nil {
			return err
		}

		err = os.Setenv("PUBSUB_PROJECT_ID", cfg.ProjectID())
		if err != nil {
			return err
		}

		err = os.Setenv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE", "1")
		if err != nil {
			return err
		}

		err = os.Setenv("FAKTORY_URL", cfg.FaktoryURL())
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

func tearDownDevelopment() runnable.Runnable {
	return runnable.Func(func(ctx context.Context) error {
		err := os.Unsetenv("STORAGE_EMULATOR_HOST")
		if err != nil {
			return err
		}

		err = os.Unsetenv("FAKTORY_PROVIDER")
		if err != nil {
			return err
		}

		err = os.Unsetenv("PUBSUB_EMULATOR_HOST")
		if err != nil {
			return err
		}

		err = os.Unsetenv("PUBSUB_PROJECT_ID")
		if err != nil {
			return err
		}

		err = os.Unsetenv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE")
		if err != nil {
			return err
		}

		return nil
	})
}
