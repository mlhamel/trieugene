package store

import (
	"context"
	"time"
)

type Store interface {
	Setup(context.Context) error
	Persist(context.Context, time.Time, string, string, interface{}) error
}
