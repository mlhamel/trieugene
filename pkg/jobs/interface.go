package jobs

import (
	"context"
	"time"
)

type Job interface {
	Name() string
	Run(ctx context.Context, args ...interface{}) error
}

type Manager interface {
	Register(Job) error
	Perform(Job, ...interface{}) error
	Run(context.Context) error
}

type Message struct {
	processedAt time.Time `json:"processed_at"`
	id          string    `json:"id"`
	data        string    `json:"data"`
}
