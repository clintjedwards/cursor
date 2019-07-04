package utils

import "errors"

// ErrPipelineNotFound is thrown when a requested entity is not found
var ErrPipelineNotFound = errors.New("pipeline not found")

// ErrPipelineExists  is thrown when an entity is found before adding another
var ErrPipelineExists = errors.New("pipeline already exists")
