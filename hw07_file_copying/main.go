package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MaksimIschenko/hw_otus_golang/hw07_file_copying/copyfile"
)

var (
	from, to      string
	limit, offset int64
)

var logger = log.New(os.Stdout, "copyfile: ", log.LstdFlags)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	logger.Println("Starting copyfile")

	err := copyfile.Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
	}

	logger.Println("Finishing copyfile")
}
