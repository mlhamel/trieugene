package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
)

type JsonStoreJob struct {
	kind  string
	cfg   *config.Config
	store store.Store
}

func NewJsonStoreJob(kind string, cfg *config.Config, store store.Store) Job {
	return &JsonStoreJob{
		kind:  kind,
		cfg:   cfg,
		store: store,
	}
}

func (r *JsonStoreJob) Kind() string {
	return r.kind
}

type dataTemp struct {
	data []interface{}
}

func (r *JsonStoreJob) Perform(ctx context.Context, args ...interface{}) error {
	r.cfg.Logger().Debug().Msgf("Running JsonStoreJob with args %v", args)

	for a := range args {
		msg, err := NewMessageFromArg(args[a])

		if err != nil {
			return err
		}

		r.cfg.Logger().Debug().Str("id", msg.ID).Int64("HappenedAt", msg.HappenedAt).Msg("Persisting data")

		datetime := time.Unix(msg.HappenedAt, 0)
		filename := fmt.Sprintf("%s/%s/%s.json", msg.Kind, datetime.Format("20060102"), datetime.Format("1504"))
		body, err := json.Marshal(msg)

		if err != nil {
			return err
		}

		bodyStr := string(body)

		if err := r.store.Persist(ctx, filename, bodyStr); err != nil {
			r.cfg.Logger().Error().Err(err).Msg("Error while trying to persist data")
			return err
		}
	}

	r.cfg.Logger().Debug().Interface("args", args).Msg("Done processing JsonStoreJob")
	return nil
}

func (r *JsonStoreJob) Run(ctx context.Context, args ...interface{}) error {
	if err := r.store.Setup(ctx); err != nil {
		r.cfg.Logger().Error().Err(err).Msg("An occured while setuping store")
		return err
	}
	if err := r.Perform(ctx, args); err != nil {
		r.cfg.Logger().Error().Err(err).Msg("An occured while running JsonStoreJob")
		return err
	}
	r.cfg.Logger().Debug().Msg("Succeed running JsonStoreJob")
	return nil
}
