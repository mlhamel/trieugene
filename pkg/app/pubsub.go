package app

import (
	"context"
	"encoding/json"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/messages"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/pior/runnable"
)

type listener struct {
	cfg *config.Config
}

func NewListener(cfg *config.Config) runnable.Runnable {
	return &listener{cfg: cfg}
}

func (l *listener) Run(ctx context.Context) error {
	pubsub, err := messages.NewPubSubListener(ctx, l.cfg)
	if err != nil {
		return err
	}

	gcs, err := store.NewGoogleCloudStorageStore(ctx, l.cfg)
	if err != nil {
		return err
	}

	err = pubsub.Listen(ctx, "rougecombien-outflow", func(ctx context.Context, msg messages.Message) error {
		data := msg.Data()
		outflow := store.Outflow{}

		err := json.Unmarshal(data, &outflow)
		if err != nil {
			return err
		}

		err = gcs.PersistOutflow(ctx, &outflow)
		if err != nil {
			return err
		}

		msg.Ack()

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
