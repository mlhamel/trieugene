package store

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/mlhamel/trieugene/pkg/config"
)

type GoogleCloudStorage struct {
	cfg    *config.Config
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGoogleCloudStorage(ctx context.Context, cfg *config.Config) (Store, error) {
	cfg.Logger().Debug().Msgf("NewGoogleCloudStorage: Initiating")
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

	cfg.Logger().Debug().Msgf("NewGoogleCloudStorage: Succeed")

	return &instance, nil
}

func (g *GoogleCloudStorage) Setup(ctx context.Context) error {
	return g.bucket.Create(ctx, g.cfg.ProjectID(), nil)
}

func (g *GoogleCloudStorage) Persist(ctx context.Context, data *Data) error {
	g.cfg.Logger().Debug().Msgf("Store/GoogleCloudStorage/Persist: Start")
	fileName := buildKey(data.Name, data.Timestamp, data.ID)
	object := g.bucket.Object(fileName)
	w := object.NewWriter(ctx)
	reader := bytes.NewReader([]byte(fmt.Sprintf("%v", data)))

	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	g.cfg.Logger().Debug().Msgf("GoogleCloudStorage/Persist: Success")

	return nil
}
