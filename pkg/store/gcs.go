package store

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"github.com/mlhamel/trieugene/pkg/config"
)

type GoogleCloudStorage struct {
	cfg    *config.Config
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGoogleCloudStorage(ctx context.Context, cfg *config.Config) (Store, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(cfg.BucketName())

	instance := GoogleCloudStorage{
		cfg:    cfg,
		client: client,
		bucket: bucket,
	}

	return &instance, nil
}

func (g *GoogleCloudStorage) Setup(ctx context.Context) error {
	return g.bucket.Create(ctx, g.cfg.ProjectID(), nil)
}

func (g *GoogleCloudStorage) Persist(ctx context.Context, timestamp time.Time, key string, data interface{}) error {
	return nil
}
