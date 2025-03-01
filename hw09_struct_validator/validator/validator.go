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

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for i, ve := range v {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ve.Field)
		sb.WriteString(": ")
		sb.WriteString(ve.Err.Error())
	}
	return sb.String()
}

func Validate(v interface{}) error {
	validationErrors := make(ValidationErrors, 0)
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		tag := fieldType.Tag.Get("validate")
		if tag == "" {
			continue
		}
		if err := ValidateField(field, tag); err != nil {
			fmt.Println(err)
			valErr := ValidationError{
				Field: fieldType.Name,
				Err:   err,
			}
			validationErrors = append(validationErrors, valErr)
		}
	}
	if len(validationErrors) > 0 {
		return fmt.Errorf("validation failed: %w", validationErrors)
	}
	return nil
}

// ValidateField validates a single field.
func ValidateField(field reflect.Value, tag string) error {
	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid rule: %s", rule)
		}
		key, value := parts[0], parts[1]
		err := validateTag(key, value, field)
		if err != nil {
			return err
		}
	}
	return nil
}

// validateTag validates a single tag.
func validateTag(key, value string, field reflect.Value) error {
	switch key {
	case "len":
		return validateLen(value, field)
	case "min":
		return validateMin(value, field)
	case "max":
		return validateMax(value, field)
	case "regexp":
		return validateRegexp(value, field)
	case "in":
		return validateIn(value, field)
	default:
		fmt.Printf("unsupported rule: %s\n", key)
		return nil
	}
}

// validateLen validates the length of the field.
func validateLen(value string, field reflect.Value) error {
	if field.Kind() != reflect.String && field.Kind() != reflect.Slice {
		return fmt.Errorf("len rule is only applicable to strings")
	}
	fieldSlice := []string{}
	if field.Kind() == reflect.Slice {
		fieldSlice = append(fieldSlice, field.Interface().([]string)...)
	} else {
		fieldSlice = append(fieldSlice, field.String())
	}

	allowedLength, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for len rule: %s", value)
	}
	for _, v := range fieldSlice {
		if len(v) != allowedLength {
			return fmt.Errorf("expected length %s, got %d", value, len(v))
		}
	}
	return nil
}

// validateMin validates the minimum value of the field.
func validateMin(value string, field reflect.Value) error {
	if field.Kind() != reflect.Int {
		return fmt.Errorf("min rule is only applicable to integers")
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for min rule: %s", value)
	}
	if field.Int() < int64(valueInt) {
		return fmt.Errorf("expected minimum %s, got %d", value, field.Int())
	}
	return nil
}

// validateMax validates the maximum value of the field.
func validateMax(value string, field reflect.Value) error {
	if field.Kind() != reflect.Int {
		return fmt.Errorf("max rule is only applicable to integers")
	}
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid value for max rule: %s", value)
	}
	if field.Int() > int64(valueInt) {
		return fmt.Errorf("expected maximum %s, got %d", value, field.Int())
	}
	return nil
}

// validateRegexp validates the field using a regular expression.
func validateRegexp(value string, field reflect.Value) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("validateRegexp: unsupported kind %s", field.Kind())
	}
	fieldValue := field.String()
	re, err := regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("invalid regular expression: %w", err)
	}
	if !re.MatchString(fieldValue) {
		return fmt.Errorf("value %q does not match pattern %q", fieldValue, value)
	}
	return nil
}

// validateIn validates the field against a list of values.
func validateIn(value string, field reflect.Value) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("validateIn: unsupported kind %s", field.Kind())
	}
	fieldValue := field.String()
	validValues := strings.Split(value, ",")
	for _, v := range validValues {
		if fieldValue == v {
			return nil
		}
	}
	return fmt.Errorf("value %q is not in allowed list %v", fieldValue, validValues)
}
