// Package currency provides ISO 4217 currency code validation.
package currency

import (
	iso4217 "golang.org/x/text/currency"
)

// IsCurrency checks if a provided string is a currency known as ISO4217 code.
func IsCurrency(text string) bool {
	_, err := iso4217.ParseISO(text)
	return err == nil
}
