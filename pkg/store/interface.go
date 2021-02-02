package store

import (
	"context"
	"time"
)

type Store interface {
	Persist(context.Context, time.Time, string, interface{}) error
}
