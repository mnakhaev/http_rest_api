package models

import validation "github.com/go-ozzo/ozzo-validation"

// requiredIf calls default validation if condition is true
func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}
		return nil
	}
}
