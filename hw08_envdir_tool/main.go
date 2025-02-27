package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/MaksimIschenko/hw_otus_golang/hw08_envdir_tool/envreader"
	"github.com/MaksimIschenko/hw_otus_golang/hw08_envdir_tool/executor"
	flag "github.com/spf13/pflag"
)

var (
	pathToEnvDir string
	cmd          []string

	wg     sync.WaitGroup
	logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
)

func main() {
	// Get positional arguments
	flag.Parse()
	positionalArgs := flag.Args()
	if len(positionalArgs) < 2 {
		logger.Error("not enough arguments")
		os.Exit(1)
	}

	// Parse arguments
	pathToEnvDir = positionalArgs[0]
	cmd = positionalArgs[1:]

	// Read environment variables
	environment, err := envreader.ReadDir(pathToEnvDir)
	if err != nil {
		logger.Error("error reading envdir", "error", err)
		os.Exit(1)
	}

	// Execute command
	wg.Add(1)
	go func() {
		returnCode := executor.RunCmd(cmd, environment)
		if returnCode != executor.OK {
			logger.Error("error executing command")
			os.Exit(1)
		}
		wg.Done()
	}()
	wg.Wait()
}
