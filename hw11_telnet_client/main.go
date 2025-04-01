package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MaksimIschenko/hw_otus_golang/hw11_telnet_client/telnet"
)

func main() {
	// parse arguments
	var timeoutStr string
	flag.StringVar(&timeoutStr, "timeout", "10s", "connection timeout")
	flag.Parse()

	// parse timeout
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse --timeout=%q: %v\n", timeoutStr, err)
		os.Exit(1)
	}

	// parse host and port
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=TIMEOUT] <host> <port>\n", os.Args[0])
		os.Exit(1)
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	// create telnet client
	client := telnet.NewTelnetClient(
		address,
		timeout,
		os.Stdin,
		os.Stdout,
	)

	// connect to server
	if err := client.Connect(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to %s: %v\n", address, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", address)

	// create context for signal handling
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	// create wait group for goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// go routine for receiving data from server
	// and writing it to stdout
	go func() {
		defer wg.Done()
		if err := client.Receive(); err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "...Error receiving: %v\n", err)
		}
		// Если err == io.EOF (или сокет закрылся), считаем, что сервер "положил трубку".
		fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
		// Отменяем контекст, чтобы завершить вторую горутину и всю программу.
		stop()
	}()

	// go routine for sending data from stdin
	go func() {
		defer wg.Done()
		if err := client.Send(); err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "...Error sending: %v\n", err)
		} else {
			// Если был EOF со стороны STDIN, сообщим и завершим.
			fmt.Fprintf(os.Stderr, "...EOF\n")
		}
		stop()
	}()

	/// wait for context cancellation
	<-ctx.Done()

	// close client connection and wait for goroutines to finish
	_ = client.Close()
	wg.Wait()
}
