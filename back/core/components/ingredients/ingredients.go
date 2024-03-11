package ingredients

import (
	"context"
	"costly/core/components/clock"
	"costly/core/components/database"
	"costly/core/components/ingredients/internal/rpst"
	"costly/core/components/logger"
	"costly/core/model"
)

type IngredientComponent interface {
	IngredientCreator
	IngredientEditor
	IngredientStockAdder
	IngredientFinder
	IngredientsFinder
}

type ingredientComponent struct {
	repository rpst.IngredientRepository
	clock      clock.Clock
}

func New(database database.Database, clock clock.Clock, logger logger.Logger) IngredientComponent {
	ingredientRepository := rpst.New(database, clock, logger)
	return &ingredientComponent{
		repository: ingredientRepository,
		clock:      clock,
	}
}

type IngredientCreator interface {
	Create(ctx context.Context, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error)
}

type CreateIngredientOptions struct {
	Name  string
	Price float64
	Unit  model.Unit
}

func (ic *ingredientComponent) Create(ctx context.Context, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error) {
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

	err := ic.repository.Add(ctx, newIngredient)

	if err != nil {
		return nil, err
	}

	return newIngredient, nil
}

type IngredientEditor interface {
	Update(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) error
}

func (ic *ingredientComponent) Update(ctx context.Context, ingredientID int64, ingredientOpts CreateIngredientOptions) error {
	err := ic.repository.Update(ctx, ingredientID, func(ingredient *model.Ingredient) error {
		ingredient.Name = ingredientOpts.Name
		ingredient.Price = ingredientOpts.Price
		ingredient.Unit = ingredientOpts.Unit
		ingredient.LastModified = ic.clock.Now()
		return nil
	})
	return err
}

type IngredientStockOptions struct {
	Units int
	Price float64
}

type IngredientStockAdder interface {
	AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error)
}

func (ic *ingredientComponent) AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error) {
	ingredientStock := &model.IngredientStock{
		ID:           -1,
		IngredientID: ingredientID,
		Units:        ingredientStockOpts.Units,
		Price:        ingredientStockOpts.Price,
		CreatedAt:    ic.clock.Now(),
	}

	err := ic.repository.AddStock(ctx, ingredientStock)

	if err != nil {
		return &model.IngredientStock{}, err
	}

	return ingredientStock, nil
}

type IngredientFinder interface {
	Find(ctx context.Context, id int64) (model.Ingredient, error)
}

func (ic *ingredientComponent) Find(ctx context.Context, id int64) (model.Ingredient, error) {
	return ic.repository.Find(ctx, id)
}

type IngredientsFinder interface {
	FindAll(ctx context.Context) ([]model.Ingredient, error)
}

func (ic *ingredientComponent) FindAll(ctx context.Context) ([]model.Ingredient, error) {
	return ic.repository.FindAll(ctx)
}
