package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	message := "Hello, OTUS!"
	reversedMessage := reverse.String(message)
	fmt.Println(reversedMessage)
}
