//go:build !solution

package requestlog

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := time.Now().Unix()
			startTime := time.Now()

			l.Info("request started",
				zap.String("path", r.URL.Path),
				zap.String("method", r.Method),
				zap.Int64("request_id", requestID),
			)

			rw := &responseWriter{w, http.StatusOK}
			defer func() {
				duration := time.Since(startTime).Seconds()

				if err := recover(); err != nil {
					l.Error("request panicked",
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
						zap.Int64("request_id", requestID),
						zap.Duration("duration", time.Duration(duration)*time.Second),
					)
					panic(r)
				}

				l.Info("request finished",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Int64("request_id", requestID),
					zap.Duration("duration", time.Duration(duration)*time.Second),
					zap.Int("status_code", rw.statusCode),
				)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
