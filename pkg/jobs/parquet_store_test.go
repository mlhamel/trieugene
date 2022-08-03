package jobs

import (
	"context"
	"path"
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

	line1 := interface{}(map[string]interface{}{
		"ProcessedAt": 20221223,
		"HappenedAt":  20221223,
		"ID":          "1",
		"Kind":        "stuff",
		"Value":       "1",
	})

	line2 := interface{}(map[string]interface{}{
		"ProcessedAt": 20221222,
		"HappennedAt": 20221222,
		"ID":          "2",
		"Kind":        "poutine",
		"Value":       "2",
	})

	err = job.Run(ctx, line1, line2)

	require.NoError(t, err)

	filename := path.Join(cfg.LocalPrefix(), "stuff", "20221223.parquet")

	require.FileExists(t, filename)
}
