package main

import (
	"encoding/json"
	"fmt"

	"github.com/MaksimIschenko/hw_otus_golang/hw09_struct_validator/validator"
)

type User struct {
	ID     string `json:"id" validate:"len:36"`
	Name   string
	Age    int             `validate:"min:18|max:50"`
	Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   UserRole        `validate:"in:admin,stuff"`
	Phones []string        `validate:"len:11"`
	meta   json.RawMessage //nolint:unused
}

type UserRole string

func main() {
	user := User{
		ID:     "123456789012345678901234567890123456",
		Name:   "John",
		Age:    25,
		Email:  "example@gmail.com",
		Role:   "admin",
		Phones: []string{"+8800553535", "+8800553536"},
	}
	err := validator.Validate(user)
	if err != nil {
		fmt.Println(err)
	}
}
