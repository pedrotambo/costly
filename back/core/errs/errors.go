package errs

import "errors"

var ErrNotFound = errors.New("entity not found")
var ErrBadOpts = errors.New("bad create entity options")
