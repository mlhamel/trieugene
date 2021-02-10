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
	ProcessedAt time.Time `json:"processed_at"`
	ID          string    `json:"id"`
	Data        string    `json:"data"`
}
