package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v *ValidationErrors) addErrors(fieldName string, errs []error) {
	for _, err := range errs {
		*v = append(*v, ValidationError{fieldName, err})
	}
}

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for i, ve := range v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("\n  ")
		sb.WriteString(ve.Field)
		sb.WriteString(": ")
		sb.WriteString(ve.Err.Error())
	}
	return sb.String()
}

func getValuesField(field reflect.Value) []any {
	values := make([]any, 0)
	if field.Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			values = append(values, convertToBaseType(field.Index(i).Interface()))
		}
	} else {
		// Для одиночного значения
		values = append(values, convertToBaseType(field.Interface()))
	}

	return values
}

func convertToBaseType(value any) any {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.String {
		return val.String()
	}
	if val.Kind() == reflect.Int {
		return int(val.Int())
	}
	return value
}

func Validate(v any) error {
	validationErrors := ValidationErrors{}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}
		values := getValuesField(field)
		for _, value := range values {
			switch v := value.(type) {
			case string:
				fieldErrors := validateStringField(v, tag)
				validationErrors.addErrors(typ.Field(i).Name, fieldErrors)
			case int:
				fieldErrors := validateIntField(v, tag)
				validationErrors.addErrors(typ.Field(i).Name, fieldErrors)
			default:
				continue
			}
		}
	}
	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors: %w", validationErrors)
	}
	return nil
}

func validateStringField(field string, tag string) []error {
	validateErrors := make([]error, 0)
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			validateErrors = append(validateErrors, fmt.Errorf("invalid rule: %s", rule))
			continue
		}
		key, value := parts[0], parts[1]
		switch key {
		case "len":
			if err := validateLen(value, field); err != nil {
				validateErrors = append(validateErrors, err)
			}
		case "regexp":
			if err := validateRegexp(value, field); err != nil {
				validateErrors = append(validateErrors, err)
			}
		case "in":
			if err := validateIn(value, field); err != nil {
				validateErrors = append(validateErrors, err)
			}
		default:
			validateErrors = append(validateErrors, fmt.Errorf("unsupported rule: %s", key))
		}
	}
	return validateErrors
}

func validateIntField(field int, tag string) []error {
	validateErrors := make([]error, 0)
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			validateErrors = append(validateErrors, fmt.Errorf("invalid rule: %s", rule))
			continue
		}
		key, value := parts[0], parts[1]
		switch key {
		case "min":
			if err := validateMin(value, field); err != nil {
				validateErrors = append(validateErrors, err)
			}
		case "max":
			if err := validateMax(value, field); err != nil {
				validateErrors = append(validateErrors, err)
			}
		case "in":
			if err := validateIn(value, field); err != nil {
				validateErrors = append(validateErrors, err)
			}
		default:
			validateErrors = append(validateErrors, fmt.Errorf("unsupported rule: %s", key))
		}
	}
	return validateErrors
}

// validateLen validates the length of the field.
func validateLen(value string, field string) error {
	allowedLength, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for len rule: %s", value)
	}
	if len(field) != allowedLength {
		return fmt.Errorf("expected length %d, got %d", allowedLength, len(field))
	}
	return nil
}

// validateMin validates the minimum value of the field.
func validateMin(value string, field int) error {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for min rule: %s", value)
	}
	if field < valueInt {
		return fmt.Errorf("allowed minimum %s, got %d", value, field)
	}
	return nil
}

// validateMax validates the maximum value of the field.
func validateMax(value string, field int) error {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for max rule: %s", value)
	}
	if field > valueInt {
		return fmt.Errorf("allowed maximum %s, got %d", value, field)
	}
	return nil
}

// validateRegexp validates the field using a regular expression.
func validateRegexp(value string, field string) error {
	re, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("invalid regular expression: %w", err)
	}
	if !re.MatchString(field) {
		return fmt.Errorf("%q invalid email", field)
	}
	return nil
}

// validateIn validates the field against a list of values.
func validateIn(value string, field interface{}) error {
	switch v := field.(type) {
	case string:
		field = v
	case int:
		field = strconv.Itoa(v)
	default:
		return fmt.Errorf("unsupported type %T", field)
	}
	validValues := strings.Split(value, ",")
	found := false
	for _, validValue := range validValues {
		if field == validValue {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%q not in allowed values %q", field, value)
	}
	return nil
}
