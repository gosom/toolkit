package scheduler

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/gosom/toolkit/pkg/errorsext"
	"github.com/gosom/toolkit/pkg/logger"
)

type task struct {
	name           string
	interval       time.Duration
	fn             TaskFunction
	runImmediately bool
}

type TaskFunction func(ctx context.Context) error

type Scheduler struct {
	log   logger.Logger
	tasks []task
}

func New(log logger.Logger) *Scheduler {
	return &Scheduler{
		log: log,
	}
}

func (s *Scheduler) AddTask(name string, interval time.Duration, fn TaskFunction, runImmediately bool) {
	s.tasks = append(s.tasks, task{
		name:           name,
		interval:       interval,
		fn:             fn,
		runImmediately: runImmediately,
	})
}

func (s *Scheduler) Run(ctx context.Context) error {
	defer func() {
		s.log.Info(ctx, "scheduler stopped")
	}()

	wg := sync.WaitGroup{}

	wg.Add(len(s.tasks))

	for _, t := range s.tasks {
		go func(t task) {
			defer func() {
				if r := recover(); r != nil {
					err := errorsext.WithStack(fmt.Errorf("panic: %v", r))
					args := []any{
						"task", t.name,
						"error", err,
					}
					s.log.Error(ctx, "panic in scheduler", args...)
					logger.ReportError(ctx, args...)
				} else {
					s.log.Info(ctx, "task stopped", "name", t.name)
				}

				wg.Done()
			}()

			execFunc := func() {
				t0 := time.Now().UTC()

				err := t.fn(ctx)

				dur := time.Now().UTC().Sub(t0)

				args := []any{
					"name", t.name,
					"duration", dur.String(),
				}
				if err != nil {
					args = append(args, "error", err)
					s.log.Error(ctx, "error running task", args...)
					logger.ReportError(ctx, args...)
				} else {
					s.log.Debug(ctx, "task finished", args...)
				}
			}

			if t.runImmediately {
				const maxSleep = 10

				randomSleep := rand.Intn(maxSleep) //nolint:gosec // this is not for security

				time.Sleep(time.Duration(randomSleep) * time.Second)

				execFunc()
			}

			ticker := time.NewTicker(t.interval)

			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					s.log.Info(ctx, "stopping task", "name", t.name)

					return
				case <-ticker.C:
					execFunc()
				}
			}
		}(t)
	}

	wg.Wait()

	return nil
}
