package domain

import (
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

type Ingredient struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Unit         Unit      `json:"unit"`
	Price        float64   `json:"price"`
	CreatedAt    time.Time `json:"created_at"`
	LastModified time.Time `json:"last_modified"`
}

type RecipeIngredient struct {
	Ingredient Ingredient `json:"ingredient"`
	Units      int        `json:"units"`
}

type Recipe struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Ingredients  []RecipeIngredient `json:"ingredients"`
	CreatedAt    time.Time          `json:"created_at"`
	LastModified time.Time          `json:"last_modified"`
}
