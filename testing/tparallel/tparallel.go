//go:build !solution

package tparallel

import "sync"

type T struct {
	parent              *T
	sequentialTest      sync.WaitGroup
	parallel            bool
	parallelTestsLaunch sync.WaitGroup
	parallelTests       sync.WaitGroup
}

func runTest(parent *T) *T {
	t := &T{
		parent:              parent,
		sequentialTest:      sync.WaitGroup{},
		parallel:            false,
		parallelTestsLaunch: sync.WaitGroup{},
		parallelTests:       sync.WaitGroup{},
	}

	t.parallelTestsLaunch.Add(1)

	return t
}

func (t *T) Parallel() {
	if t.parent == nil {
		panic("")
	}

	t.parallel = true
	t.parent.parallelTests.Add(1)
	t.parent.sequentialTest.Done()
	t.parent.parallelTestsLaunch.Wait()
}

func (t *T) Run(subtest func(t *T)) {
	t.sequentialTest.Add(1)

	go func(t *T) {
		subtest(t)

		t.parallelTestsLaunch.Done()
		t.parallelTests.Wait()

		if t.parallel {
			t.parent.parallelTests.Done()
		} else {
			t.parent.sequentialTest.Done()
		}
	}(runTest(t))

	t.sequentialTest.Wait()
}

func Run(topTests []func(t *T)) {
	t := runTest(nil)

	for _, test := range topTests {
		t.Run(test)
	}

	t.parallelTestsLaunch.Done()
	t.parallelTests.Wait()
}
