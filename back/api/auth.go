package api

import (
	"costly/core/ports/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
)

type AuthSupport struct {
	jwt *jwtauth.JWTAuth
}

func NewAuthSupport(secret []byte) *AuthSupport {
	return &AuthSupport{
		jwt: jwtauth.New("HS256", secret, nil),
	}
}

// UseMiddlewares applies the authentication support middlewares to the given router.
func (as *AuthSupport) UseMiddlewares(r chi.Router) {
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)
}

func (as *AuthSupport) PrintDebug(log logger.Logger) {
	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 1})
	log.Info("a sample jwt token", logger.Field{Key: "token", Value: tokenString})
}
