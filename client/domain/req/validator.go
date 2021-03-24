package req

import (
	"gopkg.in/go-playground/validator.v9"
)

var GValidator = validator.New()

func ValidatorStruct(s interface{}) (err error) {
	return GValidator.Struct(s)
}
