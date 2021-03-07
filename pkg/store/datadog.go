package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/mlhamel/trieugene/pkg/config"
)

type Datadog struct {
	cfg    *config.Config
	statsd *statsd.Client
}

var ErrInvalidValue = errors.Unwrap(fmt.Errorf("Invalid value for datadog"))

func NewDatadog(cfg *config.Config) Store {
	return &Datadog{cfg: cfg}
}

func (d *Datadog) Setup(ctx context.Context) error {
	statsd, err := statsd.New(d.cfg.StatsdURL(), statsd.WithNamespace("trieugene."))

	if err != nil {
		return err
	}

	d.statsd = statsd

	return nil
}

func (d *Datadog) Persist(ctx context.Context, data *Data) error {
	key := fmt.Sprintf("%s-%d", data.Name, data.Timestamp)
	valueFloat, ok := data.Value.(float64)
	if !ok {
		return ErrInvalidValue
	}
	err := d.statsd.Gauge(key, valueFloat, []string{}, 1.0)
	if err != nil {
		return err
	}

	return nil
}
