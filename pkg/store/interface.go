package store

import (
	"context"
)

type Store interface {
	Setup(ctx context.Context) error
	Persist(ctx context.Context, filename string, data string) error
}
