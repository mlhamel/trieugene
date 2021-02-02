package jobs

import (
	"context"
	"encoding/json"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
)

type OutflowJob struct {
	cfg   *config.Config
	store store.Store
}

func NewOutflowJob(cfg *config.Config, store store.Store) Job {
	return &OutflowJob{
		cfg:   cfg,
		store: store,
	}
}

func (r *OutflowJob) Name() string {
	return "outflow-rouge"
}

func (r *OutflowJob) Run(ctx context.Context, args ...interface{}) error {
	for a := range args {
		msg, ok := args[a].(Message)
		if !ok {
			return ErrInvalidMsg
		}

		var data interface{}
		err := json.Unmarshal([]byte(msg.data), &data)
		if err != nil {
			return ErrInvalidData
		}

		if err := r.store.Persist(ctx, msg.processedAt, msg.id, data); err != nil {
			return err
		}
	}

	return nil
}
