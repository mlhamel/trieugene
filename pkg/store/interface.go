package store

import (
	"context"
	"time"
)

type Store interface {
	Persist(ctx context.Context, timestamp time.Time, name string, id string, data interface{}) error
}
