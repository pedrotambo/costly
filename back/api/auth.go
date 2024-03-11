package api

import (
	"costly/core/components/logger"
	"net/http"

	"github.com/go-chi/jwtauth"
)

func NewAuthMiddleware(secret []byte, log logger.Logger) Middleware {
	h1 := jwtauth.Verifier(tokenAuth)
	h2 := jwtauth.Authenticator
	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 1})
	log.Info("a sample jwt token", logger.Field{Key: "token", Value: tokenString})
	return Middleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h1(h2(h)).ServeHTTP(w, r)
		})
	})
}
