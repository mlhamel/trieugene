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
	ProcessedAt int64       `parquet:"name=processed_at, type=INT64" json:"processed_at" csv:"processed_at"`
	HappenedAt  int64       `parquet:"name=happened_at, type=INT64" json:"happened_at" csv:"happened_at"`
	ID          string      `parquet:"name=id, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN" json:"id" csv:"id"`
	Kind        string      `parquet:"name=kind, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN" json:"kind" csv:"kind"`
	Value       interface{} `parquet:"name=value, type=BYTE_ARRAY, convertedtype=UTF8, encoding=PLAIN_DICTIONARY" json:"value" csv:"value"`
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
