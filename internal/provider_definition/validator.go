package provider_definition

import (
	"fmt"
	"reflect"
	"strings"
)

var supportedFieldTypes = map[string]struct{}{
	"string":  {},
	"integer": {},
	"number":  {},
	"boolean": {},
}

var reservedPaymentFields = map[string]struct{}{
	"id":         {},
	"provider":   {},
	"reference":  {},
	"status":     {},
	"created_at": {},
	"updated_at": {},
}

func ValidatePaymentDefinition(definition Definition) error {
	if len(definition.Payment.Fields) == 0 {
		return fmt.Errorf("provider must define at least one payment field")
	}

	for name, field := range definition.Payment.Fields {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("payment field name cannot be empty")
		}

		normalizedName := strings.ToLower(name)

		if _, exists := reservedPaymentFields[normalizedName]; exists {
			return fmt.Errorf("payment field %q is reserved by Simflux", name)
		}

		if _, exists := supportedFieldTypes[field.Type]; !exists {
			return fmt.Errorf(
				"unsupported type %q for payment field %q",
				field.Type,
				name,
			)
		}
	}

	return nil
}

func ValidatePaymentData(
	definition Definition,
	data map[string]any,
) error {
	for name, field := range definition.Payment.Fields {
		value, exists := data[name]

		if !exists {
			if field.Required {
				return fmt.Errorf("missing required payment field: %s", name)
			}
			continue
		}

		if !matchesType(value, field.Type) {
			return fmt.Errorf(
				"invalid type for payment field %q: expected %s",
				name,
				field.Type,
			)
		}
	}

	for name := range data {
		if _, exists := definition.Payment.Fields[name]; !exists {
			return fmt.Errorf("unknown payment field: %s", name)
		}
	}

	return nil
}

func matchesType(value any, expected string) bool {
	if value == nil {
		return false
	}

	kind := reflect.TypeOf(value).Kind()

	switch expected {
	case "string":
		return kind == reflect.String

	case "integer":
		return kind >= reflect.Int && kind <= reflect.Int64 ||
			kind >= reflect.Uint && kind <= reflect.Uint64

	case "number":
		return kind >= reflect.Int && kind <= reflect.Int64 ||
			kind >= reflect.Uint && kind <= reflect.Uint64 ||
			kind == reflect.Float32 || kind == reflect.Float64

	case "boolean":
		return kind == reflect.Bool
	}

	return false
}