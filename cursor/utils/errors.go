package utils

import "errors"

// ErrEntityNotFound is thrown when a requested entity is not found
var ErrEntityNotFound = errors.New("entity not found")

// ErrEntityExists  is thrown when an entity is found before adding another
var ErrEntityExists = errors.New("entity already exists")
