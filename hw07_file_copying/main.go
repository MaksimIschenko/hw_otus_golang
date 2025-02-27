package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	logger.Info("Starting copyfile")

	err := Copy(from, to, offset, limit)
	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
	}

	logger.Info("Finishing copyfile")
}
