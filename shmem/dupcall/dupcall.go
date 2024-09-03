//go:build !solution

package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	result *callResult
	mu     sync.Mutex
}

type callResult struct {
	val     any
	err     error
	done    chan struct{}
	cancel  context.CancelFunc
	counter int
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	o.mu.Lock()

	if o.result != nil {
		o.result.counter++
		o.mu.Unlock()

		return o.wait(ctx, o.result)
	}

	callCtx, cancel := context.WithCancel(context.WithoutCancel(ctx))

	c := &callResult{
		done:    make(chan struct{}),
		cancel:  cancel,
		counter: 1,
	}
	o.result = c
	o.mu.Unlock()

	go func() {
		c.val, c.err = cb(callCtx)
		close(c.done)
	}()

	return o.wait(ctx, c)
}

func (o *Call) wait(ctx context.Context, c *callResult) (v any, err error) {
	select {
	case <-c.done:
		v = c.val
		err = c.err
	case <-ctx.Done():
		err = ctx.Err()
	}

	o.mu.Lock()
	c.counter--
	if c.counter == 0 {
		c.cancel()
		o.result = nil
	}
	o.mu.Unlock()

	return v, err
}
