package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
)

func LogRequest(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("Request URL: %s", r.RequestURI))
		next.ServeHTTP(w, r)
	})
}