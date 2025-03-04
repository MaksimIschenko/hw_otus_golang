package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) { //nolint:funlen
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "valid user",
			in: User{
				ID:     "000000000000000000000000000000000000", // 36 chars
				Age:    25,
				Email:  "valid@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "invalid user ID length",
			in: User{
				ID:     "short",
				Age:    25,
				Email:  "valid@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("expected length 36, got 5")},
			},
		},
		{
			name: "invalid user age min",
			in: User{
				ID:     "000000000000000000000000000000000000",
				Age:    17,
				Email:  "valid@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("allowed minimum 18, got 17")},
			},
		},
		{
			name: "invalid user age max",
			in: User{
				ID:     "000000000000000000000000000000000000",
				Age:    51,
				Email:  "valid@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: fmt.Errorf("allowed maximum 50, got 51")},
			},
		},
		{
			name: "invalid user role",
			in: User{
				ID:     "000000000000000000000000000000000000",
				Age:    25,
				Email:  "valid@example.com",
				Role:   "user",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Role", Err: fmt.Errorf("%q not in allowed values %q", "user", "admin,stuff")},
			},
		},
		{
			name: "invalid phones length",
			in: User{
				ID:     "000000000000000000000000000000000000",
				Age:    25,
				Email:  "valid@example.com",
				Role:   "admin",
				Phones: []string{"1234567890"},
			},
			expectedErr: ValidationErrors{
				{Field: "Phones", Err: fmt.Errorf("expected length 11, got 10")},
			},
		},
		{
			name:        "valid app version",
			in:          App{Version: "1.0.0"},
			expectedErr: nil,
		},
		{
			name: "invalid app version length",
			in:   App{Version: "1.0"},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: fmt.Errorf("expected length 5, got 3")},
			},
		},
		{
			name:        "valid response code",
			in:          Response{Code: 200},
			expectedErr: nil,
		},
		{
			name: "invalid response code",
			in:   Response{Code: 400},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: fmt.Errorf("%q not in allowed values %q", "400", "200,404,500")},
			},
		},
		{
			name:        "valid token with no validation",
			in:          Token{},
			expectedErr: nil,
		},
		{
			name: "multiple errors in user",
			in: User{
				ID:     "short",
				Age:    17,
				Email:  "invalid",
				Role:   "user",
				Phones: []string{"12345"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("expected length 36, got 5")},
				{Field: "Age", Err: fmt.Errorf("allowed minimum 18, got 17")},
				{Field: "Email", Err: fmt.Errorf("%q invalid email", "invalid")},
				{Field: "Role", Err: fmt.Errorf("%q not in allowed values %q", "user", "admin,stuff")},
				{Field: "Phones", Err: fmt.Errorf("expected length 11, got 5")},
			},
		},
		{
			name: "invalid int in rule",
			in: struct {
				Number int `validate:"in:5,10,15"`
			}{7},
			expectedErr: ValidationErrors{
				{Field: "Number", Err: fmt.Errorf("%q not in allowed values %q", "7", "5,10,15")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var verrs ValidationErrors
			if !errors.As(err, &verrs) {
				t.Fatalf("expected ValidationErrors, got %T", err)
			}

			var expectedErrs ValidationErrors
			if !errors.As(tt.expectedErr, &expectedErrs) {
				t.Fatalf("invalid test setup: expectedErr must be ValidationErrors")
			}

			if len(verrs) != len(expectedErrs) {
				t.Fatalf("expected %d errors, got %d: %v", len(expectedErrs), len(verrs), verrs)
			}

			for i, expected := range expectedErrs {
				actual := verrs[i]
				if actual.Field != expected.Field || actual.Err.Error() != expected.Err.Error() {
					t.Errorf("error %d:\nexpected: %s: %v\ngot:      %s: %v",
						i, expected.Field, expected.Err, actual.Field, actual.Err)
				}
			}
		})
	}
}

func TestRuntimeError(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr string
	}{
		{
			name: "unsupported rule",
			in: struct {
				Field string `validate:"unknown:value"`
			}{"test"},
			expectedErr: "unsupported rule: unknown",
		},
		{
			name: "invalid rule format",
			in: struct {
				Field string `validate:"len|max:5"`
			}{"test"},
			expectedErr: "invalid rule: len",
		},
		{
			name: "invalid min value (non-integer)",
			in: struct {
				Age int `validate:"min:abc"`
			}{20},
			expectedErr: "strconv.Atoi: parsing \"abc\": invalid syntax",
		},
		{
			name: "invalid regexp",
			in: struct {
				Field string `validate:"regexp:*invalid"`
			}{"test"},
			expectedErr: "error parsing regexp: missing argument to repetition operator: `*`",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error %q, got %q", tt.expectedErr, err.Error())
			}
		})
	}
}
