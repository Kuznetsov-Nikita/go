//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type KeyLock struct {
	mu    sync.Mutex
	locks map[string]chan struct{}
}

func New() *KeyLock {
	return &KeyLock{locks: make(map[string]chan struct{})}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	sortedKeys := make([]string, len(keys))
	copy(sortedKeys, keys)
	sort.Strings(sortedKeys)

	acquired := make([]chan struct{}, 0)

	canceled = false
	unlock = func() {
		l.mu.Lock()
		for key := range l.locks {
			delete(l.locks, key)
		}
		l.mu.Unlock()

		for _, lock := range acquired {
			close(lock)
		}
	}

	for _, key := range sortedKeys {
		for {
			l.mu.Lock()
			otherLock, isLocked := l.locks[key]
			lock := make(chan struct{})

			if !isLocked {
				l.locks[key] = lock
			}
			l.mu.Unlock()

			if !isLocked {
				acquired = append(acquired, lock)
				break
			}

			select {
			case <-cancel:
				unlock()
				canceled = true
				return
			case <-otherLock:
				continue
			}
		}
	}

	return
}
