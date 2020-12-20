package store

import (
	"context"
	"testing"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNewGoogleCloudStorageStore(t *testing.T) {
	store, err := NewGoogleCloudStorageStore(context.Background(), config.NewConfig())

	require.NoError(t, err)
	require.NotNil(t, store)
}

func TestGoogleCloudStorageStoreSetup(t *testing.T) {
	ctx := context.Background()
	store, _ := NewGoogleCloudStorageStore(ctx, config.NewConfig())

	err := store.Setup(ctx)

	require.NoError(t, err)
}
