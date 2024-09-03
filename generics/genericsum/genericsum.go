//go:build !solution

package genericsum

import (
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a <= b {
		return a
	} else {
		return b
	}
}

func SortSlice[T constraints.Ordered](a []T) {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
}

func MapsEqual[K comparable, V comparable](a, b map[K]V) bool {
	if len(a) != len(b) {
		return false
	}

	for key, value := range a {
		secondValue, ok := b[key]
		if !ok {
			return false
		}

		if secondValue != value {
			return false
		}
	}

	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, elem := range s {
		if v == elem {
			return true
		}
	}

	return false
}

func MergeChans[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	wg.Add(len(chs))
	for _, c := range chs {
		go func(c <-chan T) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

type Number interface {
	constraints.Integer | constraints.Complex
}

func conjugate(num any) any {
	switch v := num.(type) {
	case int:
		return v
	case int8:
		return v
	case int16:
		return v
	case int32:
		return v
	case int64:
		return v
	case uint:
		return v
	case uint8:
		return v
	case uint16:
		return v
	case uint32:
		return v
	case uint64:
		return v
	case float32:
		return v
	case float64:
		return v
	case complex64:
		return complex(real(v), -imag(v))
	case complex128:
		return complex(real(v), -imag(v))
	default:
		panic("Unsupported type")
	}
}

func IsHermitianMatrix[T Number](m [][]T) bool {
	n := len(m)
	if n == 0 || n != len(m[0]) {
		return false
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if m[i][j] != conjugate(m[j][i]) {
				return false
			}
		}
	}

	return true
}
