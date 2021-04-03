package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/mitchellh/mapstructure"
	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
)

type CsvStoreJob struct {
	kind  string
	cfg   *config.Config
	store store.Store
}

func NewCsvStoreJob(kind string, cfg *config.Config, store store.Store) Job {
	return &CsvStoreJob{
		kind:  kind,
		cfg:   cfg,
		store: store,
	}
}

func (c *CsvStoreJob) Kind() string {
	return c.kind
}

func (c *CsvStoreJob) Perform(ctx context.Context, args ...interface{}) error {
	c.cfg.Logger().Debug().Str("job", "CsvStoreJob").Msgf("Running with %d messages", len(args))
	messages := []*Message{}

	for a := range args {

		c.cfg.Logger().Debug().Str("job", "CsvStoreJob").Msg("Parsing arguments")
		data, ok := args[a].([]interface{})
		if !ok {
			c.cfg.Logger().Error().Err(ErrInvalidMsg).Msg("Error while parsing args")
			return ErrInvalidMsg
		}
		c.cfg.Logger().Debug().Str("job", "CsvStoreJob").Int("length", len(data)).Msg("Succeed: Parsing arguments")

		for b := range data {
			var msg Message
			err := mapstructure.Decode(data[b], &msg)
			if err != nil {
				c.cfg.Logger().Error().Str("job", "CsvStoreJob").Err(err).Msg("Error while decoding msg")
				return err
			}

			messages = append(messages, &msg)
		}
	}

	var first = messages[0]

	datetime := time.Unix(first.HappenedAt, 0)
	filename := fmt.Sprintf("%s/%s.csv", first.Kind, datetime.Format("20060102"))
	body, err := gocsv.MarshalString(&messages)

	if err != nil {
		return fmt.Errorf("Error while marshaling into a csv: %w", err)
	}

	c.cfg.Logger().Debug().Str("filename", filename).Msg("Persisting csv")
	if err := c.store.Persist(ctx, filename, body); err != nil {
		c.cfg.Logger().Error().Err(ErrInvalidData).Msg("Failed at persisting csv")
		return ErrInvalidData
	}

	return nil
}

func (c *CsvStoreJob) Run(ctx context.Context, args ...interface{}) error {
	if err := c.store.Setup(ctx); err != nil {
		c.cfg.Logger().Error().Err(err).Msg("An occured while setuping store")
		return err
	}
	if err := c.Perform(ctx, args); err != nil {
		c.cfg.Logger().Error().Err(err).Msg("An occured while running CsvStoreJob")
		return err
	}
	c.cfg.Logger().Debug().Msg("Succeed running CsvStoreJob")
	return nil
}
