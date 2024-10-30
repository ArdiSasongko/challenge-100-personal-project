package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMessages []string
		for _, errV := range validationErrors {
			switch errV.Tag() {
			case "required":
				errMessages = append(errMessages, fmt.Sprintf("%s is required", errV.Field()))
			case "min":
				errMessages = append(errMessages, fmt.Sprintf("%s is min", errV.Field()))
			case "max":
				errMessages = append(errMessages, fmt.Sprintf("%s is max", errV.Field()))
			case "number":
				errMessages = append(errMessages, fmt.Sprintf("%s is number", errV.Field()))
			case "alphanum":
				errMessages = append(errMessages, fmt.Sprintf("%s is alphanum", errV.Field()))
			default:
				errMessages = append(errMessages, fmt.Sprintf("%s is not valid", errV.Field()))
			}
		}
		return errors.New(strings.Join(errMessages, ","))
	}
	return nil
}
