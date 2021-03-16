package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mlhamel/trieugene/pkg/config"
	"github.com/mlhamel/trieugene/pkg/store"
	"github.com/pior/runnable"
)

type Store struct {
	cfg   *config.Config
	store store.Store
	kind  string
	key   string
	value string
}

type data struct {
	Timestamp int64       `json:"timestamp"`
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
}

func NewStore(cfg *config.Config, store store.Store, kind string, key string, value string) runnable.Runnable {
	return &Store{
		cfg:   cfg,
		store: store,
		kind:  kind,
		key:   key,
		value: value,
	}
}

func (s *Store) Run(ctx context.Context) error {
	s.cfg.Logger().Debug().Msgf("Apps/Store/Run: Start")

	var result interface{}

	err := json.Unmarshal([]byte(s.value), &result)
	if err != nil {
		return fmt.Errorf("Error while unmarshaling value (%s): %w", s.value, err)
	}

	timestamp := time.Now()
	filename := fmt.Sprintf("%s/%s/%s.json", s.kind, timestamp.Format("20060102"), timestamp.Format("1504"))
	dataStr, err := json.Marshal(data{
		Timestamp: time.Now().Unix(),
		Name:      s.kind,
		ID:        s.key,
		Value:     result,
	})

	if err != nil {
		return fmt.Errorf("Error while marshaling data: %w", err)
	}

	err = s.store.Persist(context.Background(), filename, string(dataStr))
	if err != nil {
		return err
	}

	s.cfg.Logger().Debug().Msgf("Apps/Store/Run: Success")

	return nil
}
