package main

import (
	"context"
	"flag"
	"github.com/johnewart/freighter/server"
	"sync"
	"zombiezen.com/go/log"
)

func main() {
	wg := sync.WaitGroup{}

	rootDir := "/Users/johnewart/.cache/registry/"
	port := flag.Int("port", 50051, "The server port")
	flag.Parse()
	ctx := context.Background()

	fs := server.NewFreighterServer(rootDir)
	wg.Add(2)
	go func() {
		if err := fs.Serve("0.0.0.0", *port); err != nil {
			log.Errorf(ctx, "Error: %v", err)
			wg.Done()
		}
	}()

	go func() {
		if err := fs.ServeRegistry("0.0.0.0", 1338); err != nil {
			log.Errorf(ctx, "Failed to start registry: %v", err)
			wg.Done()
		}
	}()

	log.Infof(ctx, "Freighter is up and running...")
	wg.Wait()
}
