package messages

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/mlhamel/trieugene/pkg/config"
)

type PubSubListener struct {
	cfg    *config.Config
	client *pubsub.Client
}

type PubSubMessage struct {
	original *pubsub.Message
}

func NewPubSubMessage(original *pubsub.Message) *PubSubMessage {
	return &PubSubMessage{original: original}
}

func (p *PubSubMessage) Ack() {
	p.original.Ack()
}

func (p *PubSubMessage) Nack() {
	p.original.Nack()
}

func (p *PubSubMessage) Data() []byte {
	return p.original.Data
}

func NewPubSubListener(ctx context.Context, cfg *config.Config) (Listener, error) {
	client, err := pubsub.NewClient(ctx, cfg.ProjectID())

	if err != nil {
		return nil, err
	}

	return &PubSubListener{
		cfg:    cfg,
		client: client,
	}, nil
}

func (p *PubSubListener) Listen(ctx context.Context, topic string, callback ListenerFuncType) error {
	var sub *pubsub.Subscription

	sub, err := p.find(ctx, topic)

	if err != nil {
		return fmt.Errorf("could not search for subscription: %v", err)
	}

	if sub == nil {
		sub, err = p.create(ctx, topic)

		if err != nil {
			return fmt.Errorf("could not create subscription: %v", err)
		}
	}

	cctx, cancel := context.WithCancel(ctx)
	return sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		callback(cctx, NewPubSubMessage(msg))
		defer cancel()
	})
}

func (p *PubSubListener) find(ctx context.Context, topic string) (*pubsub.Subscription, error) {
	it := p.client.Subscriptions(ctx)

	for {
		subscription, err := it.Next()

		if err != nil {
			return nil, nil
		}

		if subscription.ID() == topic {
			return subscription, nil
		}
	}

	return nil, nil
}

func (p *PubSubListener) create(ctx context.Context, topic string) (*pubsub.Subscription, error) {
	t := p.client.Topic(topic)
	exists, err := t.Exists(ctx)

	if err != nil {
		return nil, err
	}

	if !exists {
		if t, err = p.client.CreateTopic(ctx, topic); err != nil {
			return nil, err
		}
	}

	sub, err := p.client.CreateSubscription(ctx, topic, pubsub.SubscriptionConfig{
		Topic:       t,
		AckDeadline: 10 * time.Second,
	})

	return sub, err
}
