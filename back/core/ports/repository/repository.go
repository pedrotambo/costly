package repository

import (
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"errors"
)

type Repository struct {
	IngredientRepository
	RecipeRepository
}

var ErrNotFound = errors.New("entity not found")
var ErrBadOpts = errors.New("bad create entity options")

func New(db *database.Database, clock clock.Clock) *Repository {
	return &Repository{
		NewIngredientRepository(db, clock),
		NewRecipeRepository(db, clock),
	}
}
