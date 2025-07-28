package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func Init() {
	Validate = validator.New()
	_ = Validate.RegisterValidation("monthyear", validateMonthYear)
}

func validateMonthYear(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^(0[1-9]|1[0-2])-20\d{2}$`)
	return re.MatchString(fl.Field().String())
}
