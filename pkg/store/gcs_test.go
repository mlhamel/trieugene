package store

import (
	"context"
	"testing"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNewGoogleCloudStorage(t *testing.T) {
	store, err := NewGoogleCloudStorage(context.Background(), config.NewConfig())

	require.NoError(t, err)
	require.NotNil(t, store)
}

func TestGoogleCloudStorageSetup(t *testing.T) {
	ctx := context.Background()

	store, err := NewGoogleCloudStorage(ctx, config.NewConfig())
	require.NoError(t, err)

	err = store.Setup(ctx)
	require.NoError(t, err)
}
