package jobs

import (
	"context"
)

type Job interface {
	Run(ctx context.Context, args ...interface{}) error
}

type Manager interface {
	Register(string, Job) error
	Perform(string, Job, *Message) error
	Run(context.Context) error
}

type Message struct {
	ProcessedAt int64  `json:"processed_at" mapstructure:"processed_at"`
	HappenedAt  int64  `json:"happened_at" mapstructure:"happened_at"`
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Data        string `json:"data"`
}
