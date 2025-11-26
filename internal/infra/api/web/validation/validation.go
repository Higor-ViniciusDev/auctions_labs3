package validation

import (
	"encoding/json"
	"errors"

	"github.com/Higor-ViniciusDev/auction_labs3/configuration/rest_err"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	validator_en "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate   = validator.New()
	Translator ut.Translator
)

func init() {
	if value, ok := binding.Validator.Engine().(*validator.Validate); ok {
		En := en.New()
		uni := ut.New(En, En)
		Translator, _ = uni.GetTranslator("en")
		_ = validator_en.RegisterDefaultTranslations(value, Translator)
	}
}

func ValidateErr(err error) *rest_err.RestErr {
	// Check for JSON unmarshal type errors and validation errors
	var jsonErr *json.UnmarshalTypeError
	// Check for validation errors
	var jsonValidation validator.ValidationErrors

	if errors.As(err, &jsonErr) {
		return rest_err.NewNotFoundError("Invalid type error")
	} else if errors.As(err, &jsonValidation) {
		errorCauses := []rest_err.Causes{}

		for _, eThr := range err.(validator.ValidationErrors) {
			errorCause := rest_err.Causes{
				Field:   eThr.Field(),
				Message: eThr.Translate(Translator),
			}
			errorCauses = append(errorCauses, errorCause)
		}

		return rest_err.NewBadRequestError("Validation error", errorCauses...)
	} else {
		return rest_err.NewBadRequestError("Invalid request body")
	}
}
