package api

import (
	"costly/api/handlers"

	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/repository"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
}

type Middleware func(http.Handler) http.Handler

func NewRouter(database *database.Database, clock clock.Clock, repository *repository.Repository, as *AuthSupport, middlewares ...Middleware) chi.Router {
	r := chi.NewRouter()

	for _, m := range middlewares {
		r.Use(m)
	}

	// cors middleware for development
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	// TODO: change this to use zerolog
	r.Use(middleware.Logger)

	// authenticated routes
	r.Group(func(r chi.Router) {
		as.UseMiddlewares(r)
		// ingredients
		r.Get("/ingredients", handlers.GetIngredientsHandler(repository))
		r.Post("/ingredients", handlers.CreateIngredientHandler(repository))
		r.Get("/ingredients/{ingredientID}", handlers.GetIngredientHandler(repository))

		// recipes
		r.Post("/recipes", handlers.CreateRecipeHandler(repository))
		r.Get("/recipes", handlers.GetRecipesHandler(repository))
		r.Get("/recipes/{recipeID}", handlers.GetRecipeHandler(repository))
	})

	return r
}
