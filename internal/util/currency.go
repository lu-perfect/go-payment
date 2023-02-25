package util

import "github.com/go-playground/validator/v10"

const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, RUB:
		return true
	}
	return false
}

var CurrencyValidator validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return IsSupportedCurrency(currency)
	}
	return false
}
