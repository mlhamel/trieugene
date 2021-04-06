package jobs

import "context"

type ExpectsJob struct {
	Called bool
}

func (e *ExpectsJob) Kind() string {
	return "expected-job"
}

func (e *ExpectsJob) Perform(ctx context.Context, args ...interface{}) error {
	e.Called = true
	return nil
}

func (e *ExpectsJob) Run(ctx context.Context, args ...interface{}) error {
	return e.Perform(ctx, args)
}
