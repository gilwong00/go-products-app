package data

import "github.com/go-playground/validator"

type ValidationError struct {
	validator.FieldError
}
