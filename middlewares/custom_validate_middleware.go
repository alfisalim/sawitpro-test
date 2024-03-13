package middlewares

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type ValidationHandler struct {
	Tag          string
	Func         validator.Func
	ErrorMessage string
}

type CustomValidatorInterface interface {
	Validate(interface{}) error
}

// CustomValidator is struct used to create request validator
type CustomValidator struct {
	Validator           *validator.Validate
	ValidatorHandlerMap map[string]ValidationHandler
}

// NewValidator is function to create custom validator struct
func NewValidator() *CustomValidator {
	return &CustomValidator{
		Validator: validator.New(),
	}
}

// Validate is function to validate all request based on struct
// you can see how to using it on
// https://github.com/go-playground/validator
func (c *CustomValidator) Validate(i interface{}) error {
	customValidatePassword(c)
	if err := c.Validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range validationErrors {
				var message string
				if handler, ok := c.ValidatorHandlerMap[err.Tag()]; ok {
					message = fmt.Sprintf("invalid fields '%s' with message: %s", err.Field(), handler.ErrorMessage)
				} else {
					message = buildMessageWithTag(err)
				}

				return errors.New(message)
			}
		}
	}

	return nil
}

func buildMessageWithTag(v validator.FieldError) string {
	if v.Param() != "" {
		switch v.Tag() {
		case "min":
			return fmt.Sprintf("invalid field '%s' with issue minimum char is %s", v.Field(), v.Param())
		case "max":
			return fmt.Sprintf("invalid field '%s' with issue maximum char is %s", v.Field(), v.Param())
		case "startswith":
			return fmt.Sprintf("invalid field '%s' with issue must be starts with %s", v.Field(), v.Param())
		case "oneof":
			return fmt.Sprintf("invalid field '%s' with issue at least contains %s", v.Field(), v.Param())
		default:
			return fmt.Sprintf("%s", v.Error())
		}
	} else {
		switch v.Tag() {
		case "required":
			return fmt.Sprintf("required field '%s'", v.Field())
		case "numeric":
			return fmt.Sprintf("invalid field '%s', must be numeric", v.Field())
		case "validpasswd":
			return fmt.Sprintf("invalid field '%s', please use combination of alphanumeric and special character with lowercase and uppercase", v.Field())
		default:
			return fmt.Sprintf("%s", v.Error())
		}
	}
}

func customValidatePassword(c *CustomValidator) {
	passwordString := `^(.*[a-z])(.*[A-Z])(.*\d)(.*[@$!%*?&;])` // regex that compiles
	passwordRegex := regexp.MustCompile(passwordString)

	c.Validator.RegisterValidation("validpasswd", func(fl validator.FieldLevel) bool {
		return passwordRegex.MatchString(fl.Field().String())
	})
}
