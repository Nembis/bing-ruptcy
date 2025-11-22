package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type attributeLogger string

const enrichedLoggerKey attributeLogger = "logger.attribute.enrich"

type responseRecorder struct {
	http.ResponseWriter
	Status int
	Bytes  int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		Status:         200,
	}
}

func (r *responseRecorder) writeHeader(code int) {
	r.Status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.Bytes += n
	return n, err
}

func midddlewareLoggerEnricher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := newResponseRecorder(w)

		attrs := []slog.Attr{
			slog.String("request_method", r.Method),
			slog.String("request_path", r.URL.Path),
			slog.String("query_string", r.URL.RawQuery),
			slog.String("request_scheme", r.URL.Scheme),
			slog.String("request_host", r.Host),
			slog.String("request_protocol", r.Proto),
			slog.String("request_content_type", r.Header.Get("Content-Type")),
			slog.Int64("request_content_length", r.ContentLength),
			slog.String("remote_ip", r.RemoteAddr),
		}

		anyAttr := make([]any, 0, len(attrs))
		for _, attr := range attrs {
			anyAttr = append(anyAttr, attr)
		}

		enrichedLogger := baseLogger.With(anyAttr...)
		ctx := context.WithValue(r.Context(), enrichedLoggerKey, enrichedLogger)
		next.ServeHTTP(rec, r.WithContext(ctx))

		enrichedLogger.InfoContext(ctx, "Request Completed",
			"StatusCode", rec.Status,
			"Duration", time.Since(start).Milliseconds(),
		)
	})
}

func getEnrichedLogger(ctx context.Context) *slog.Logger {
	if enrichedLogger, ok := ctx.Value(enrichedLoggerKey).(*slog.Logger); ok {
		return enrichedLogger
	}
	return baseLogger
}
