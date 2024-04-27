package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

type loggedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func LogRequest(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		wrapped := &loggedWriter{
			ResponseWriter: w,
			statusCode: http.StatusOK,
		}
		next.ServeHTTP(wrapped, r)
		logger.Info("REQUEST", "Status Code", fmt.Sprint(wrapped.statusCode), "Method", r.Method, "Route", r.URL.Path)
	})
}