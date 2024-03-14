package api

import (
	"costly/core/ports/logger"
	"net/http"
)

func NewLoggerMiddleware(logger logger.Logger) Middleware {
	return Middleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithContext(r.Context())
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}
