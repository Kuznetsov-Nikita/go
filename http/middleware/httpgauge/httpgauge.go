//go:build !solution

package httpgauge

import (
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Gauge struct {
	metrics map[string]int
	mu      sync.Mutex
}

func New() *Gauge {
	return &Gauge{metrics: make(map[string]int)}
}

func (g *Gauge) Snapshot() map[string]int {
	return g.metrics
}

// ServeHTTP returns accumulated statistics in text format ordered by pattern.
//
// For example:
//
//	/a 10
//	/b 5
//	/c/{id} 7
func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	patterns := []string{}

	for pattern := range g.metrics {
		patterns = append(patterns, pattern)
	}

	sort.Strings(patterns)

	for _, pattern := range patterns {
		_, _ = w.Write([]byte(pattern + " " + strconv.Itoa(g.metrics[pattern]) + "\n"))
	}
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rctx := chi.RouteContext(r.Context())
			pattern := rctx.RoutePattern()

			g.mu.Lock()
			g.metrics[pattern]++
			g.mu.Unlock()
		}()
		next.ServeHTTP(w, r)
	})
}
