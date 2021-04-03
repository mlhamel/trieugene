package jobs

import (
	"context"

	"github.com/mitchellh/mapstructure"
)

type Job interface {
	Kind() string
	Run(ctx context.Context, args ...interface{}) error
}

type Manager interface {
	Register(Job) error
	Perform(Job, ...*Message) error
	Run(context.Context) error
}

type Message struct {
	ProcessedAt int64       `json:"processed_at" mapstructure:"processed_at" csv:"processed_at"`
	HappenedAt  int64       `json:"happened_at" mapstructure:"happened_at" csv:"happened_at"`
	ID          string      `json:"id" csv:"id"`
	Kind        string      `json:"kind" csv:"kind"`
	Value       interface{} `json:"value" csv:"value"`
}

func NewMessageFromArg(arg interface{}) (*Message, error) {
	var msg Message

	data, ok := arg.([]interface{})
	if !ok {
		return nil, ErrInvalidMsg
	}

	raw, ok := data[0].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidMsg
	}

	err := mapstructure.Decode(raw, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
