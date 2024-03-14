package model

import (
	"costly/core/errs"
	"time"
)

type Unit string

const (
	Gram       Unit = "gr"
	Kilogram   Unit = "kg"
	Liter      Unit = "L"
	Milliliter Unit = "ml"
	Units      Unit = "units"
)

type ID int64

type Ingredient struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Unit         Unit      `json:"unit"`
	Price        float64   `json:"price"`
	UnitsInStock int       `json:"units_in_stock"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}

func NewIngredient(name string, unit Unit, price float64, now time.Time) (*Ingredient, error) {
	if name == "" {
		return &Ingredient{}, errs.ErrBadName
	}
	if unit != "gr" {
		return &Ingredient{}, errs.ErrBadUnit
	}

	if price <= 0 {
		return &Ingredient{}, errs.ErrBadPrice
	}
	return &Ingredient{
		ID:           -1,
		Name:         name,
		Unit:         unit,
		Price:        price,
		UnitsInStock: 0,
		CreatedAt:    now,
		LastModified: now,
	}, nil
}

type IngredientStock struct {
	ID           int64     `json:"id"`
	IngredientID int64     `json:"ingredient_id"`
	Units        int       `json:"units"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
}

func NewIngredientStock(ingredientID int64, units int, price float64, now time.Time) *IngredientStock {
	return &IngredientStock{
		ID:           -1,
		IngredientID: ingredientID,
		Units:        units,
		Price:        price,
		CreatedAt:    now,
	}
}

type RecipeIngredientView struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Units int     `json:"units"`
}

type RecipeView struct {
	ID           int64                  `json:"id"`
	Name         string                 `json:"name"`
	Ingredients  []RecipeIngredientView `json:"ingredients"`
	CreatedAt    time.Time              `json:"created_at"`
	LastModified time.Time              `json:"last_modified"`
}

func (recipe *RecipeView) Cost() float64 {
	cost := 0.0

	for _, ingredient := range recipe.Ingredients {
		cost += ingredient.Price * float64(ingredient.Units)
	}

	return cost
}

type RecipeSales struct {
	ID        int64
	RecipeID  int64
	Units     int
	CreatedAt time.Time
}

func NewRecipeSales(recipeID int64, units int, now time.Time) *RecipeSales {
	return &RecipeSales{
		ID:        -1,
		RecipeID:  recipeID,
		Units:     units,
		CreatedAt: now,
	}
}

type RecipeIngredient struct {
	ID    int64 `json:"id"`
	Units int   `json:"units"`
}

type Recipe struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Ingredients  []RecipeIngredient `json:"ingredients"`
	CreatedAt    time.Time          `json:"created_at"`
	LastModified time.Time          `json:"last_modified"`
}

func NewRecipe(name string, ingredients []RecipeIngredient, now time.Time) (*Recipe, error) {
	if name == "" {
		return &Recipe{}, errs.ErrBadName
	}
	if len(ingredients) == 0 {
		return &Recipe{}, errs.ErrBadIngrs
	}
	return &Recipe{
		ID:           -1,
		Name:         name,
		Ingredients:  ingredients,
		CreatedAt:    now,
		LastModified: now,
	}, nil
}
