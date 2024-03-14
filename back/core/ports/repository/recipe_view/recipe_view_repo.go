package recipeviewrepo

import (
	"context"
	"costly/core/model"
	"costly/core/ports/database"
	"time"
)

type RecipeViewRepository interface {
	FindIngredients(ctx context.Context, recipeID int64) ([]model.RecipeIngredientView, error)
	FindAll(ctx context.Context) ([]model.RecipeView, error)
}

type repository struct {
	db database.Database
}

func New(db database.Database) RecipeViewRepository {
	return &repository{db}
}

func (r *repository) FindIngredients(ctx context.Context, recipeID int64) ([]model.RecipeIngredientView, error) {
	recipeIngredients, err := database.QueryAndMap(ctx, r.db, mapToRecipeIngredientView, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", recipeID)
	if err != nil {
		return nil, err
	}
	return recipeIngredients, nil
}

func (r *repository) FindAll(ctx context.Context) ([]model.RecipeView, error) {
	recipesDB, err := database.QueryAndMap(ctx, r.db, mapToRecipeDB, "SELECT * FROM recipe")
	if err != nil {
		return nil, err
	}
	recipes := []model.RecipeView{}
	for _, recipeDB := range recipesDB {
		recipeIngredients, err := database.QueryAndMap(ctx, r.db, mapToRecipeIngredientView, "SELECT i.*, ri.units FROM ingredient i JOIN recipe_ingredient ri ON i.id = ri.ingredient_id AND ri.recipe_id = ?", recipeDB.id)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, model.RecipeView{
			ID:           recipeDB.id,
			Name:         recipeDB.name,
			Ingredients:  recipeIngredients,
			CreatedAt:    recipeDB.createdAt,
			LastModified: recipeDB.lastModified,
		})
	}
	return recipes, nil
}

func mapToRecipeIngredientView(rowScanner database.RowScanner) (model.RecipeIngredientView, error) {
	var ingredient model.Ingredient
	var recipeUnits int
	err := rowScanner.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Unit, &ingredient.Price, &ingredient.CreatedAt, &ingredient.LastModified, &ingredient.UnitsInStock, &recipeUnits)
	return model.RecipeIngredientView{
		ID:    ingredient.ID,
		Name:  ingredient.Name,
		Price: ingredient.Price,
		Units: recipeUnits,
	}, err
}

type recipeDB struct {
	id           int64
	name         string
	createdAt    time.Time
	lastModified time.Time
}

type recipeViewDB struct {
	id           int64
	name         string
	createdAt    time.Time
	lastModified time.Time
}

func mapToRecipeDB(rowScanner database.RowScanner) (recipeDB, error) {
	var recipe recipeDB
	err := rowScanner.Scan(&recipe.id, &recipe.name, &recipe.createdAt, &recipe.lastModified)
	return recipe, err
}
