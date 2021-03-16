package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
)

type StoreJob struct {
	kind  string
	cfg   *config.Config
	store store.Store
}

func NewStoreJob(kind string, cfg *config.Config, store store.Store) Job {
	return &StoreJob{
		kind:  kind,
		cfg:   cfg,
		store: store,
	}
}

func (r *StoreJob) Kind() string {
	return r.kind
}

type dataTemp struct {
	data []interface{}
}

func (r *StoreJob) Perform(ctx context.Context, args ...interface{}) error {
	r.cfg.Logger().Debug().Msgf("Running StoreJob with args %v", args)

	for a := range args {
		var msg Message
		data, ok := args[a].([]interface{})
		if !ok {
			r.cfg.Logger().Error().Err(ErrInvalidMsg).Msg("Invalid message")
			return ErrInvalidMsg
		}

		raw, ok := data[0].(map[string]interface{})
		if !ok {
			r.cfg.Logger().Error().Err(ErrInvalidMsg).Msg("Invalid message")
			return ErrInvalidMsg
		}

		r.cfg.Logger().Debug().Interface("Raw", raw).Msg("Decoding arguments")
		err := mapstructure.Decode(raw, &msg)
		if err != nil {
			r.cfg.Logger().Error().Err(err).Msg("Error while trying to decode arguments")
			return err
		}

		r.cfg.Logger().Debug().Str("id", msg.ID).Int64("HappenedAt", msg.HappenedAt).Msg("Persisting data")

		datetime := time.Unix(msg.HappenedAt, 0)
		filename := fmt.Sprintf("%s/%s/%s.json", msg.Kind, datetime.Format("20060102"), datetime.Format("1504"))
		body, err := json.Marshal(data)

		if err != nil {
			return err
		}

		bodyStr := string(body)

		if err := r.store.Persist(ctx, filename, bodyStr); err != nil {
			r.cfg.Logger().Error().Err(err).Msg("Error while trying to persist data")
			return err
		}
	}

	r.cfg.Logger().Debug().Interface("args", args).Msg("Done processing StoreJob")
	return nil
}

func (r *StoreJob) Run(ctx context.Context, args ...interface{}) error {
	if err := r.store.Setup(ctx); err != nil {
		r.cfg.Logger().Error().Err(err).Msg("An occured while setuping store")
		return err
	}
	if err := r.Perform(ctx, args); err != nil {
		r.cfg.Logger().Error().Err(err).Msg("An occured while running StoreJob")
		return err
	}
	r.cfg.Logger().Debug().Msg("Succeed running StoreJob")
	return nil
}
