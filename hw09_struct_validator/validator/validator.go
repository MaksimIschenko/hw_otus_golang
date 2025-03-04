package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Validation error type.
type ValidationError struct {
	Field string
	Err   error
}

// ValidationErrors is a list of validation errors.
type ValidationErrors []ValidationError

// Add errors to the validation errors.
func (v *ValidationErrors) addErrors(fieldName string, errs []error) {
	for _, err := range errs {
		*v = append(*v, ValidationError{fieldName, err})
	}
}

// Error returns the validation errors as a string.
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

// Get the values of the field.
func getValuesField(field reflect.Value) []any {
	values := make([]any, 0)
	if field.Kind() == reflect.Slice {
		for i := 0; i < field.Len(); i++ {
			values = append(values, convertToBaseType(field.Index(i).Interface()))
		}
	} else {
		values = append(values, convertToBaseType(field.Interface()))
	}

	return values
}

// Convert the value to the base type.
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

// Validate validates the struct fields.
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
				fieldErrors, runtimeError := validateStringField(v, tag)
				if runtimeError != nil {
					return runtimeError
				}
				validationErrors.addErrors(typ.Field(i).Name, fieldErrors)
			case int:
				fieldErrors, runtimeError := validateIntField(v, tag)
				if runtimeError != nil {
					return runtimeError
				}
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

// validateStringField validates the string field.
func validateStringField(field string, tag string) ([]error, error) {
	rules := strings.Split(tag, "|")
	validateErrors := make([]error, 0)
	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			return []error{}, fmt.Errorf("invalid rule: %s", rule)
		}
		key, value := parts[0], parts[1]
		switch key {
		case "len":
			validErr, runtimeErr := validateLen(value, field)
			if runtimeErr != nil {
				return []error{}, runtimeErr
			}
			if validErr != nil {
				validateErrors = append(validateErrors, validErr)
			}
		case "regexp":
			validErr, runtimeErr := validateRegexp(value, field)
			if runtimeErr != nil {
				return []error{}, runtimeErr
			}
			if validErr != nil {
				validateErrors = append(validateErrors, validErr)
			}
		case "in":
			validErr, runtimeErr := validateIn(value, field)
			if runtimeErr != nil {
				return []error{}, runtimeErr
			}
			if validErr != nil {
				validateErrors = append(validateErrors, validErr)
			}
		default:
			return []error{}, fmt.Errorf("unsupported rule: %s", key)
		}
	}
	return validateErrors, nil
}

// validateIntField validates the integer field.
func validateIntField(field int, tag string) ([]error, error) {
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
			validErr, runtimeErr := validateMin(value, field)
			if runtimeErr != nil {
				return []error{}, runtimeErr
			}
			if validErr != nil {
				validateErrors = append(validateErrors, validErr)
			}

		case "max":
			validErr, runtimeErr := validateMax(value, field)
			if runtimeErr != nil {
				return []error{}, runtimeErr
			}
			if validErr != nil {
				validateErrors = append(validateErrors, validErr)
			}
		case "in":
			validErr, runtimeErr := validateIn(value, field)
			if runtimeErr != nil {
				return []error{}, runtimeErr
			}
			if validErr != nil {
				validateErrors = append(validateErrors, validErr)
			}
		default:
			return []error{}, fmt.Errorf("unsupported rule: %s", key)
		}
	}
	return validateErrors, nil
}

// validateLen validates the length of the field.
func validateLen(value string, field string) (validErr error, runtimeErr error) {
	allowedLength, err := strconv.Atoi(value)
	if err != nil {
		runtimeErr = fmt.Errorf("invalid value for len rule: %s", value)
		return
	}
	if len(field) != allowedLength {
		validErr = fmt.Errorf("expected length %d, got %d", allowedLength, len(field))
	}
	return
}

// validateMin validates the minimum value of the field.
func validateMin(value string, field int) (validErr error, runtimeErr error) {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		runtimeErr = err
		return
	}
	if field < valueInt {
		validErr = fmt.Errorf("allowed minimum %s, got %d", value, field)
	}
	return
}

// validateMax validates the maximum value of the field.
func validateMax(value string, field int) (validErr error, runtimeErr error) {
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		runtimeErr = err
		return
	}
	if field > valueInt {
		validErr = fmt.Errorf("allowed maximum %s, got %d", value, field)
	}
	return
}

// validateRegexp validates the field using a regular expression.
func validateRegexp(value string, field string) (validErr error, runtimeErr error) {
	re, err := regexp.Compile(value)
	if err != nil {
		runtimeErr = err
		return
	}
	if !re.MatchString(field) {
		validErr = fmt.Errorf("%q invalid email", field)
	}
	return
}

// validateIn validates the field against a list of allowed values.
func validateIn(value string, field interface{}) (validErr error, runtimeErr error) {
	switch v := field.(type) {
	case string:
		field = v
	case int:
		field = strconv.Itoa(v)
	default:
		runtimeErr = fmt.Errorf("unsupported type %T", field)
		return
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
		validErr = fmt.Errorf("%q not in allowed values %q", field, value)
	}
	return
}
