package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/libmonsoon-dev/gomut/src/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	errCh := make(chan error, 1)
	go func() { errCh <- wg.Wait() }()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	application, err := app.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	wg.Go(func() error { return application.Run(ctx) })

	fmt.Printf("Got signal: %v\n", <-signals)

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Print(err)
		}
	case <-time.After(5 * time.Second):
		panic("Timeout")
	}

}
