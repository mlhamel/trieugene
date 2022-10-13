package jobs

import (
	"context"
	"path"
	"testing"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/stretchr/testify/require"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func lines() []interface{} {
	return append(make([]interface{}, 0),
		interface{}(map[string]interface{}{
			"ProcessedAt": 20221222,
			"HappennedAt": 20221222,
			"ID":          "2",
			"Kind":        "poutine",
			"Value":       "2",
		}),
		interface{}(map[string]interface{}{
			"ProcessedAt": 20221223,
			"HappenedAt":  20221223,
			"ID":          "1",
			"Kind":        "stuff",
			"Value":       "1",
		}),
	)
}

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

	err = job.Run(ctx, lines()...)

	require.NoError(t, err)

	filename := path.Join(cfg.LocalPrefix(), "stuff", "20221223.parquet")

	require.FileExists(t, filename)
}

func TestParquetStoreContent(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig()
	st := store.NewLocal(cfg)

	err := st.Setup(ctx)

	require.NoError(t, err)

	job := NewParquetStoreJob("stuff", cfg, st)

	err = job.Run(ctx, lines()...)

	require.NoError(t, err)

	filename := path.Join(cfg.LocalPrefix(), "stuff", "20221223.parquet")

	fr, err := local.NewLocalFileReader(filename)
	require.NoError(t, err)

	pr, err := reader.NewParquetReader(fr, new(Message), 2)
	require.NoError(t, err)

	require.Equal(t, 2, int(pr.GetNumRows()))
}
