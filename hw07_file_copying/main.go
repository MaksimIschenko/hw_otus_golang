package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	fmt.Printf("from: %v\n", from)
	fmt.Printf("to: %v\n", to)
	fmt.Printf("limit: %v\n", limit)
	fmt.Printf("offset: %v\n", offset)
}
