package errs

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("entity not found")
var ErrBadOpts = errors.New("")

func newBadOptsError(msg string) error {
	return fmt.Errorf("%w%s", ErrBadOpts, msg)
}

var ErrBadName = newBadOptsError("name is invalid")
var ErrBadUnit = newBadOptsError("unit is invalid")
var ErrBadPrice = newBadOptsError("price is invalid")
var ErrBadIngrs = newBadOptsError("recipe must have at least one ingredient")
var ErrBadStockUnits = newBadOptsError("units should be more than 0")
