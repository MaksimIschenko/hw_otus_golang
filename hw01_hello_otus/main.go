package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	message := "Hello, OTUS!"
	fmt.Println(reverse.String(message))
}
