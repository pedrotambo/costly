package ingredients

import (
	"context"
	"costly/core/components/ingredients/internal/rpst"
	"costly/core/errs"
	"costly/core/model"
	"costly/core/ports/clock"
	"costly/core/ports/database"
	"costly/core/ports/logger"
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
	ingredientRepository := rpst.New(database, logger)
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

func (opts CreateIngredientOptions) validate() error {
	if opts.Name == "" {
		return errs.ErrBadName
	}

	if opts.Unit != "gr" {
		return errs.ErrBadUnit
	}

	if opts.Price <= 0 {
		return errs.ErrBadPrice
	}
	return nil
}

func (ic *ingredientComponent) Create(ctx context.Context, ingredientOpts CreateIngredientOptions) (*model.Ingredient, error) {
	if err := ingredientOpts.validate(); err != nil {
		return &model.Ingredient{}, err
	}
	now := ic.clock.Now()
	newIngredient := model.NewIngredient(ingredientOpts.Name, ingredientOpts.Unit, ingredientOpts.Price, now)
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
	if err := ingredientOpts.validate(); err != nil {
		return err
	}
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

func (opts IngredientStockOptions) validate() error {
	if opts.Units <= 0 {
		return errs.ErrBadStockUnits
	}

	if opts.Price <= 0 {
		return errs.ErrBadPrice
	}
	return nil
}

type IngredientStockAdder interface {
	AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error)
}

func (ic *ingredientComponent) AddStock(ctx context.Context, ingredientID int64, ingredientStockOpts IngredientStockOptions) (*model.IngredientStock, error) {
	if err := ingredientStockOpts.validate(); err != nil {
		return &model.IngredientStock{}, err
	}
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
