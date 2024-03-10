package usecases

import (
	"context"
	"costly/core/domain"
	"costly/core/ports/clock"
	"costly/core/ports/rpst"
)

type CreateIngredientOptions struct {
	Name  string
	Price float64
	Unit  domain.Unit
}

type IngredientUseCases interface {
	IngredientCreator
	IngredientEditor
}

type IngredientCreator interface {
	CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (*domain.Ingredient, error)
}

type IngredientEditor interface {
	EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (*domain.Ingredient, error)
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

func (ic *ingredientUseCases) CreateIngredient(ctx context.Context, ingredientOpts CreateIngredientOptions) (*domain.Ingredient, error) {
	now := ic.clock.Now()
	newIngredient := &domain.Ingredient{
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

func (ic *ingredientUseCases) EditIngredient(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) (*domain.Ingredient, error) {
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
