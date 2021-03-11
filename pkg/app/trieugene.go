package app

import (
	"context"
	"os"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/pior/runnable"
)

type trieugene struct {
	cfg *config.Config
}

type trieugeneDev struct {
	cfg *config.Config
}

type trieugeneStore struct {
	cfg   *config.Config
	kind  string
	key   string
	value string
}

func NewTrieugene(cfg *config.Config) runnable.Runnable {
	return &trieugene{cfg: cfg}
}

func NewTrieugeneDev(cfg *config.Config) runnable.Runnable {
	return &trieugeneDev{cfg: cfg}
}

func NewTrieugeneStore(cfg *config.Config, kind string, key string, value string) runnable.Runnable {
	return &trieugeneStore{
		cfg:   cfg,
		kind:  kind,
		key:   key,
		value: value,
	}
}

func (t *trieugene) Run(ctx context.Context) error {
	store := store.NewS3(&store.S3Params{
		AccessKey: t.cfg.GCSAccessKey(),
		SecretKey: t.cfg.GCSAccessSecret(),
		URL:       t.cfg.GCSURL(),
		Bucket:    t.cfg.S3Bucket(),
		Region:    t.cfg.S3Region(),
	})

	err := NewFaktory(t.cfg, store).Run(ctx)

	return err
}

func (t *trieugeneDev) Run(ctx context.Context) error {
	store := store.NewS3(&store.S3Params{
		AccessKey:        t.cfg.S3AccessKey(),
		SecretKey:        t.cfg.S3SecretKey(),
		URL:              t.cfg.S3URL(),
		Bucket:           t.cfg.S3Bucket(),
		Region:           t.cfg.S3Region(),
		DisableSSL:       true,
		S3ForcePathStyle: true,
	})

	run(setupDevelopment(t.cfg))

	err := NewFaktory(t.cfg, store).Run(ctx)

	run(tearDownDevelopment())

	return err
}

func (t *trieugeneStore) Run(ctx context.Context) error {
	store := store.NewS3(&store.S3Params{
		AccessKey:        t.cfg.S3AccessKey(),
		SecretKey:        t.cfg.S3SecretKey(),
		URL:              t.cfg.S3URL(),
		Bucket:           t.cfg.S3Bucket(),
		Region:           t.cfg.S3Region(),
		DisableSSL:       true,
		S3ForcePathStyle: true,
	})

	run(setupDevelopment(t.cfg))

	err := NewStore(t.cfg, store, t.kind, t.key, t.value).Run(ctx)

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
