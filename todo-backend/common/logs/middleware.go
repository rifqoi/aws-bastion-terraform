package logs

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func NewMiddleware(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// Record ResponseWriter
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				reqLogger := logger.With(
					zap.String("proto", r.Proto),
					zap.Int64("ts", t1.Unix()),
					zap.String("method", r.Method),
					zap.Duration("lat", time.Since(t1)),
					zap.String("reqID", middleware.GetReqID(r.Context())),
					zap.String("path", r.URL.Path),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)

				reqLogger.Info("SERVED")

			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
