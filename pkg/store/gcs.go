package store

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/mlhamel/trieugene/pkg/config"
)

type GoogleCloudStorageStore struct {
	cfg    *config.Config
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGoogleCloudStorageStore(ctx context.Context, cfg *config.Config) (Store, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	bucket := client.Bucket(cfg.BucketName())

	instance := GoogleCloudStorageStore{
		cfg:    cfg,
		client: client,
		bucket: bucket,
	}

	return &instance, nil
}

func (g *GoogleCloudStorageStore) Setup(ctx context.Context) error {
	return g.bucket.Create(ctx, g.cfg.ProjectID(), nil)
}

func (g *GoogleCloudStorageStore) PersistOutflow(ctx context.Context, outflow *Outflow) error {
	return nil
}
