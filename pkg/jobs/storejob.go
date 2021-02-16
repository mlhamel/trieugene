package jobs

import (
	"context"
	"encoding/json"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
)

type StoreJob struct {
	cfg   *config.Config
	store store.Store
}

func NewStoreJob(cfg *config.Config, store store.Store) Job {
	return &StoreJob{
		cfg:   cfg,
		store: store,
	}
}

func (r *StoreJob) Perform(ctx context.Context, args ...interface{}) error {
	r.cfg.Logger().Debug().Msgf("Running StoreJob with args %v", args)

	for a := range args {
		var msg = Message{}
		jsonString, err := json.Marshal(args[a])
		if err != nil {
			r.cfg.Logger().Error().Msgf("Invalid message '%s': %w", args[a], err)
			return ErrInvalidMsg
		}

		json.Unmarshal(jsonString, &msg)

		r.cfg.Logger().Debug().Msg("Unmarshaling data for persistence")
		var data interface{}
		err = json.Unmarshal([]byte(msg.Data), &data)
		if err != nil {
			return ErrInvalidData
		}

		r.cfg.Logger().Debug().Msgf("Persisting data for %d processedAt %s", msg.ID, msg.ProcessedAt.String())
		if err := r.store.Persist(ctx, msg.ProcessedAt, msg.Kind, msg.ID, data); err != nil {
			return err
		}
	}

	r.cfg.Logger().Debug().Msgf("Done processing StoreJob with args %v", args)
	return nil
}

func (r *StoreJob) Run(ctx context.Context, args ...interface{}) error {
	if err := r.Perform(ctx, args); err != nil {
		r.cfg.Logger().Error().Err(err).Msg("An occured while running StoreJob")
		return err
	}
	return nil
}
