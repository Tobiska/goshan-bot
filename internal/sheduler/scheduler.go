package sheduler

import (
	"context"
	"time"
)

type Scheduler struct {
	period time.Duration
}

func New(period time.Duration) *Scheduler {
	return &Scheduler{
		period: period,
	}
}

func (s *Scheduler) Run(ctx context.Context, runnable func(ctx context.Context)) error {
	for ctx.Err() == nil {
		runnable(ctx)
		s.wait()
	}

	return nil
}

func (s *Scheduler) wait() {
	time.Sleep(s.period)
}
