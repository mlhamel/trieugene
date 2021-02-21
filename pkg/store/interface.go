package store

import (
	"context"
)

type Store interface {
	Setup(context.Context) error
	Persist(context.Context, int64, string, string, interface{}) error
}
