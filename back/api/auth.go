package api

import (
	"costly/core/ports/logger"
	"net/http"

	"github.com/go-chi/jwtauth"
)

type AuthSupport struct {
	jwt        *jwtauth.JWTAuth
	middleware Middleware
}

func NewAuthSupport(secret []byte) *AuthSupport {
	h1 := jwtauth.Verifier(tokenAuth)
	h2 := jwtauth.Authenticator
	return &AuthSupport{
		jwt: jwtauth.New("HS256", secret, nil),
		middleware: Middleware(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				h1(h2(h)).ServeHTTP(w, r)
			})
		}),
	}
}

func (as *AuthSupport) PrintDebug(log logger.Logger) {
	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 1})
	log.Info("a sample jwt token", logger.Field{Key: "token", Value: tokenString})
}
