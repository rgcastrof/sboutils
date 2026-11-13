package pkg

import (
	"sync"
	"os"
	"fmt"
	"context"
)

type semaphore struct {
	channel chan struct{}
}

func (s *semaphore) Acquire(ctx context.Context) error {
	select {
		case<-ctx.Done():
			return ctx.Err()
		case s.channel <- struct{}{}:
			return nil
	}
}

func (s *semaphore) Release() {
	<-s.channel
}

type Dispatcher[T any] struct {
	items []T
	sem *semaphore
	wg sync.WaitGroup
	mu sync.Mutex
	firstErr error
}

func NewDispatcher[T any](values []T, concurrencyLimit int) *Dispatcher[T] {
	if concurrencyLimit <= 0 {
		fmt.Fprintf(os.Stderr, "ERROR: concurrencyLimit must be  greater than 0\n")
		os.Exit(1)
	}

	return &Dispatcher[T]{
		items: values,
		sem: &semaphore{channel: make(chan struct{}, concurrencyLimit)},
	}
}

func (dispatcher *Dispatcher[T]) Run(ctx context.Context, op func(param T) error) error {
	ctx, cancel := context.WithCancel(ctx)
	var acquireErr error
	defer cancel()

	for _, item := range dispatcher.items {
		if err := dispatcher.sem.Acquire(ctx); err != nil {
			acquireErr = err
			break
		}

		dispatcher.wg.Add(1)
		go func(value T) {
			defer dispatcher.wg.Done()
			defer dispatcher.sem.Release()

			if err := op(value); err != nil {
				dispatcher.mu.Lock()
				if dispatcher.firstErr == nil {
					dispatcher.firstErr = err
					cancel()
				}
				dispatcher.mu.Unlock()
			}
		}(item)
	}

	dispatcher.wg.Wait()

	dispatcher.mu.Lock()
	defer dispatcher.mu.Unlock()
	if dispatcher.firstErr != nil {
		return dispatcher.firstErr
	}

	if acquireErr != nil {
		return acquireErr
	}

	return nil
}
