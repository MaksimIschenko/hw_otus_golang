package main

import (
	"fmt"

	"github.com/beaconsoftwarellc/gadget/stringutil" //nolint
)

func main() {
	message := "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(message))
}
