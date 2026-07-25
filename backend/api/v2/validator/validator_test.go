package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

type TestDecimalStruct struct {
	Val        decimal.Decimal  `validate:"decimal"`
	PtrVal     *decimal.Decimal `validate:"omitempty,decimal"`
	StringVal  string           `validate:"decimal"`
	InvalidVal string           `validate:"decimal"`
}

func TestDecimalValidator(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("decimal", Decimal)

	val1 := decimal.NewFromFloat(123.45)
	validStruct := TestDecimalStruct{
		Val:        val1,
		PtrVal:     &val1,
		StringVal:  "123.45",
		InvalidVal: "12.34",
	}

	err := v.Struct(validStruct)
	if err != nil {
		t.Errorf("expected validStruct to pass validation, got: %v", err)
	}

	invalidStruct := TestDecimalStruct{
		Val:        val1,
		PtrVal:     nil,
		StringVal:  "123.45",
		InvalidVal: "not_a_decimal",
	}

	err = v.Struct(invalidStruct)
	if err == nil {
		t.Errorf("expected invalidStruct to fail validation on 'not_a_decimal', got nil")
	}
}
