package store

import (
	"context"
)

type Data struct {
	Timestamp int64       `json:"timestamp"`
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Value     interface{} `json:"value"`
}

type Store interface {
	Setup(ctx context.Context) error
	Persist(context.Context, *Data) error
}
