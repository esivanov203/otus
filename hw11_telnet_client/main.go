package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	//nolint:depguard
	"github.com/spf13/pflag"
)

func main() {
	ts := pflag.String("timeout", "10s", "timeout in seconds")
	pflag.Parse()

	args := pflag.Args()
	if len(args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [--timeout=10s] host port\n", os.Args[0])
		os.Exit(1)
	}

	h := args[0]
	p := args[1]
	t, err := time.ParseDuration(*ts)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Invalid timeout: %s\nUsage: --timeout=10s\n", *ts)
		os.Exit(1)
	}

	tc := NewTelnetClient(fmt.Sprintf("%s:%s", h, p), t, os.Stdin, os.Stdout)
	err = tc.Connect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed connect to %s %s: %v", h, p, err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		_, _ = fmt.Fprintf(os.Stderr, "Нажато CTRL+C\n")
		_ = tc.Close()
		os.Exit(0)
	}()

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := tc.Send()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Отправка завершена: %v\n", err)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Отправлено EOF\n")
		}
		_ = tc.Close()
	}()

	go func() {
		defer wg.Done()
		err := tc.Receive()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Соединение закрыто сервером: %v\n", err)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Соединение закрыто сервером\n")
		}
		_ = tc.Close()
	}()

	wg.Wait()
}
