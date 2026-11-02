package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return []Result[R]{}, nil
	}

	n := len(in)
	results := make([]Result[R], n)

	w := workers
	if w > n {
		w = n
	}

	jobs := make(chan int)
	var wg sync.WaitGroup
	wg.Add(w)

	for i := 0; i < w; i++ {
		go func() {
			defer wg.Done()
			for idx := range jobs {
				if ctx.Err() != nil {
					return
				}
				v, err := fn(ctx, in[idx])
				if err != nil {
					var zero R
					v = zero
				}
				results[idx] = Result[R]{Value: v, Err: err}
			}
		}()
	}

	func() {
		defer close(jobs)
		for i := 0; i < n; i++ {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
