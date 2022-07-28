package jobs

import (
	"context"
	"testing"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestParquetStoreWithoutData(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	st := store.NewLocal(cfg)

	err := st.Setup(ctx)

	require.NoError(t, err)

	job := NewParquetStoreJob("stuff", cfg, st)

	require.Equal(t, "stuff", job.Kind())

	err = job.Run(ctx)

	require.NoError(t, err)
}

func TestParquetStoreWithtData(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	st := store.NewLocal(cfg)

	err := st.Setup(ctx)

	require.NoError(t, err)

	job := NewParquetStoreJob("stuff", cfg, st)

	require.Equal(t, "stuff", job.Kind())

	line := interface{}(map[string]string{
		"ProcessedAt": "20221223",
		"HappennedAt": "20221223",
		"ID":          "1",
		"Kind":        "stuff",
		"Value":       "1",
	})

	err = job.Run(ctx, line)

	require.NoError(t, err)
}
