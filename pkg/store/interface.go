package store

import "context"

type Store interface {
	Setup(context.Context) error
	PersistOutflow(context.Context, *Outflow) error
}
