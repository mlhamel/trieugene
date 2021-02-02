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
	key   string
	value string
}

func NewStore(cfg *config.Config, key string, value string) runnable.Runnable {
	store := store.NewS3(cfg)

	return &Store{cfg: cfg, store: store, key: key, value: value}
}

func (s *Store) Run(ctx context.Context) error {
	var result interface{}

	err := json.Unmarshal([]byte(s.value), &result)
	if err != nil {
		return fmt.Errorf("Error while unmarshaling value (%s): %w", s.value, err)
	}

	return s.store.Persist(ctx, time.Now(), s.key, result)
}
