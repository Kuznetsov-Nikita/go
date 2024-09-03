//go:build !solution

package waitgroup

type generation struct {
	wait chan struct{}
	n    int
}

func newGeneration() generation {
	return generation{wait: make(chan struct{})}
}

func (g generation) end() {
	close(g.wait)
}

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
type WaitGroup chan generation

// New creates WaitGroup.
func New() *WaitGroup {
	wg := make(WaitGroup, 1)
	g := newGeneration()
	g.end()
	wg <- g
	return &wg
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	g := <-*wg
	if g.n == 0 {
		g = newGeneration()
	}
	g.n += delta
	if g.n < 0 {
		panic("negative WaitGroup counter")
	}
	if g.n == 0 {
		g.end()
	}
	*wg <- g
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	g := <-*wg
	wait := g.wait
	*wg <- g
	<-wait
}
