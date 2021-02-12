package jobs

import (
	"context"
	"time"
)

type Job interface {
	Run(ctx context.Context, args ...interface{}) error
}

type Manager interface {
	Register(string, Job) error
	Perform(string, Job, ...interface{}) error
	Run(context.Context) error
}

type Message struct {
	ProcessedAt time.Time `json:"processed_at"`
	ID          string    `json:"id"`
	Data        string    `json:"data"`
}
