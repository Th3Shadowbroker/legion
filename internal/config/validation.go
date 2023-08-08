package config

import (
	validation "github.com/go-playground/validator/v10"
)

var validator *validation.Validate

func init() {
	validator = validation.New()
}

func ValidatePreset(preset *Preset) validation.ValidationErrors {
	var err = validator.Struct(preset)
	return err.(validation.ValidationErrors)
}
