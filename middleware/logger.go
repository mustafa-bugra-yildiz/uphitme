package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggerResponseWriter(w http.ResponseWriter) *loggerResponseWriter {
	return &loggerResponseWriter{w, http.StatusOK}
}

func (lrw *loggerResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := newLoggerResponseWriter(w)

		start := time.Now()
		handler.ServeHTTP(lrw, r)
		duration := time.Since(start)

		log.Println(
			r.Method,
			r.URL.String()+":",
			lrw.statusCode,
			http.StatusText(lrw.statusCode)+",",
			duration,
		)
	})
}
