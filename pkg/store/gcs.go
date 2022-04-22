package store

import (
	"bytes"
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/mlhamel/trieugene/pkg/config"
)

type GoogleCloudStorage struct {
	cfg        *config.Config
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
}

func NewGoogleCloudStorage(ctx context.Context, cfg *config.Config) (Store, error) {
	cfg.Logger().Debug().Msgf("NewGoogleCloudStorage: Initiating")

	instance := GoogleCloudStorage{
		bucketName: cfg.BucketName(),
		cfg:        cfg,
		client:     nil,
		bucket:     nil,
	}

	cfg.Logger().Debug().Msgf("NewGoogleCloudStorage: Succeed")

	return &instance, nil
}

func (g *GoogleCloudStorage) Initialized() bool {
	return g.client != nil && g.bucket != nil
}

func (g *GoogleCloudStorage) Init(ctx context.Context) error {
	g.cfg.Logger().Debug().Msgf("Store/GoogleCloudStorage/Init: Start")

	if g.Initialized() {
		return nil
	}
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	g.client = client
	g.bucket = g.client.Bucket(g.bucketName)

	g.cfg.Logger().Debug().Msgf("GoogleCloudStorage/Init: Success")

	return nil
}

func (g *GoogleCloudStorage) Setup(ctx context.Context) error {
	g.cfg.Logger().Debug().Msgf("Store/GoogleCloudStorage/Setup: Start")
	err := g.bucket.Create(ctx, g.cfg.ProjectID(), nil)
	if err != nil {
		return err
	}
	g.cfg.Logger().Debug().Msgf("Store/GoogleCloudStorage/Setup: Start")
	return nil
}

func (g *GoogleCloudStorage) Persist(ctx context.Context, filename string, data string) error {
	g.cfg.Logger().Debug().Msgf("Store/GoogleCloudStorage/Persist: Start")
	g.Init(ctx)

	object := g.bucket.Object(filename)
	w := object.NewWriter(ctx)
	reader := bytes.NewReader([]byte(data))

	if _, err := io.Copy(w, reader); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	g.cfg.Logger().Debug().Msgf("GoogleCloudStorage/Persist: Success")

	return nil
}
