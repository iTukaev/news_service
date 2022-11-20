package customerrors

import "github.com/pkg/errors"

var (
	ErrNewsNotFound      = errors.New("news not found")
	ErrNewsAlreadyExists = errors.New("news already exists")
	ErrTimeout           = errors.New("deadline exceeded")
	ErrUnexpected        = errors.New("unexpected error")
)
