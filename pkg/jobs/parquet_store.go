package jobs

import (
	"bytes"
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/xitongsys/parquet-go/writer"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
)

type ParquetStoreJob struct {
	kind  string
	cfg   *config.Config
	store store.Store
}

func NewParquetStoreJob(kind string, cfg *config.Config, store store.Store) Job {
	return &ParquetStoreJob{
		kind:  kind,
		cfg:   cfg,
		store: store,
	}
}

func (p *ParquetStoreJob) Kind() string {
	return p.kind
}

func (c *ParquetStoreJob) Perform(ctx context.Context, args ...interface{}) error {
	c.cfg.Logger().Debug().Str("job", "ParquetStoreJob").Msgf("Performing with %d messages", len(args))
	var first *Message

	bf := bytes.NewBufferString("")
	writer, err := writer.NewParquetWriterFromWriter(bf, &Message{}, 5)

	if err != nil {
		return err
	}

	for a := range args {
		data, ok := args[a].([]interface{})

		if !ok {
			c.cfg.Logger().Error().Err(ErrInvalidMsg).Msg("Error while parsing args")
			return ErrInvalidMsg
		}

		for b := range data {
			var msg Message
			err := mapstructure.Decode(data[b], &msg)

			if err != nil {
				c.cfg.Logger().Error().Str("job", "ParquetStoreJob").Err(err).Msg("Error while decoding msg")
				return err
			}

			if err = writer.Write(msg); err != nil {
				return err
			}

			if first == nil {
				first = &msg
			}
		}
	}

	if err = writer.WriteStop(); err != nil {
		return err
	}

	if first == nil {
		return nil
	}

	filename := fmt.Sprintf("%s/%d.parquet", first.Kind, first.HappenedAt)
	bytes := bf.Bytes()
	body := string(bytes)

	if err := c.store.Persist(ctx, filename, body); err != nil {
		c.cfg.Logger().Error().Err(ErrInvalidData).Msg("Failed at persisting parquet")
		return ErrInvalidData
	}

	return nil
}

func (p *ParquetStoreJob) Run(ctx context.Context, args ...interface{}) error {
	p.cfg.Logger().Debug().Str("job", "ParquetStoreJob").Msgf("Running with %d messages", len(args))

	if err := p.store.Setup(ctx); err != nil {
		p.cfg.Logger().Error().Err(err).Msg("An occured while setuping store")
		return err
	}

	if err := p.Perform(ctx, args); err != nil {
		p.cfg.Logger().Error().Err(err).Msg("An occured while running ParquetStoreJob")
		return err
	}
	p.cfg.Logger().Debug().Msg("Succeed running ParquetStoreJob")

	return nil
}
