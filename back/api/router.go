package api

import (
	"costly/api/handlers"

	"costly/core/usecases"
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

func NewRouter(useCases usecases.UseCases, authMiddleware Middleware, middlewares ...Middleware) http.Handler {
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
		r.Use(authMiddleware)
		// ingredients
		r.Get("/ingredients", handlers.GetIngredientsHandler(useCases))
		r.Post("/ingredients", handlers.CreateIngredientHandler(useCases))
		r.Get("/ingredients/{ingredientID}", handlers.GetIngredientHandler(useCases))
		r.Put("/ingredients/{ingredientID}", handlers.EditIngredientHandler(useCases))
		r.Post("/ingredients/{ingredientID}/stock", handlers.AddIngredientStockHandler(useCases))

		// recipes
		r.Post("/recipes", handlers.CreateRecipeHandler(useCases))
		r.Get("/recipes", handlers.GetRecipesHandler(useCases))
		r.Get("/recipes/{recipeID}", handlers.GetRecipeHandler(useCases))
	})

	return r
}
