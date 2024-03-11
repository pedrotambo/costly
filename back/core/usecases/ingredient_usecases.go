package usecases

import (
	"context"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/rpst"
)

type IngredientUseCases interface {
	IngredientCreator
	IngredientEditor
	IngredientStockAdder
}

type IngredientCreator interface {
	CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error)
}

type IngredientEditor interface {
	EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error)
}

type IngredientStockAdder interface {
	AddIngredientStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error)
}

type ingredientUseCases struct {
	repository rpst.IngredientRepository
	clock      clock.Clock
}

func NewIngredientUseCases(repository rpst.IngredientRepository, clock clock.Clock) IngredientUseCases {
	return &ingredientUseCases{
		repository: repository,
		clock:      clock,
	}
}

type CreateIngredientOptions struct {
	Name  string
	Price float64
	Unit  model.Unit
}

func (ic *ingredientUseCases) CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error) {
	now := ic.clock.Now()
	newIngredient := &model.Ingredient{
		ID:           -1,
		Name:         ingredientOpts.Name,
		Unit:         ingredientOpts.Unit,
		Price:        ingredientOpts.Price,
		UnitsInStock: 0,
		CreatedAt:    now,
		LastModified: now,
	}

	err := ic.repository.SaveIngredient(ctx, newIngredient)

	if err != nil {
		return nil, err
	}

	return newIngredient, nil
}

func (ic *ingredientUseCases) EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error) {
	ingredient, err := ic.repository.GetIngredient(ctx, ingredientID)
	if err != nil {
		return &ingredient, err
	}
	ingredient.Name = ingredientOpts.Name
	ingredient.Price = ingredientOpts.Price
	ingredient.Unit = ingredientOpts.Unit
	ingredient.LastModified = ic.clock.Now()

	err = ic.repository.UpdateIngredient(ctx, &ingredient)

	if err != nil {
		return &ingredient, err
	}

	return &ingredient, nil
}

type IngredientStockOptions struct {
	Units int
	Price float64
}

func (ic *ingredientUseCases) AddIngredientStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error) {
	ingredientStock := &model.IngredientStock{
		ID:           -1,
		IngredientID: ingredientID,
		Units:        ingredientStockOpts.Units,
		Price:        ingredientStockOpts.Price,
		CreatedAt:    ic.clock.Now(),
	}

	err := ic.repository.SaveIngredientStock(ctx, ingredientStock)

	if err != nil {
		return &model.IngredientStock{}, err
	}

	return ingredientStock, nil
}
