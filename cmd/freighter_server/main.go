package main

import (
	"context"
	"flag"
	"sync"

	"github.com/johnewart/freighter/server"
	"zombiezen.com/go/log"
)

var (
	cacheRoot    = flag.String("cacheroot", "./freightercache", "the root of the image layer store")
	grpcPort     = flag.Int("grpcport", 50051, "The FS gRPC server port")
	registryPort = flag.Int("registryport", 1338, "The registry gRPC server port")
)

func main() {
	wg := sync.WaitGroup{}

	flag.Parse()

	if cacheRoot == nil || *cacheRoot == "" {
		log.Errorf(context.Background(), "Please specify a cache root directory with -cacheroot")
		return
	}

	ctx := context.Background()

	fs := server.NewFreighterServer(*cacheRoot)
	wg.Add(2)
	go func() {
		if err := fs.Serve("0.0.0.0", *grpcPort); err != nil {
			log.Errorf(ctx, "Error: %v", err)
			wg.Done()
		}
	}()

	go func() {
		if err := fs.ServeRegistry("0.0.0.0", *registryPort); err != nil {
			log.Errorf(ctx, "Failed to start registry: %v", err)
			wg.Done()
		}
	}()

	log.Infof(ctx, "Freighter is up and running!")
	log.Infof(ctx, "gRPC server listening on port %d", *grpcPort)
	log.Infof(ctx, "Registry server listening on port %d", *registryPort)
	log.Infof(ctx, "Cache root: %s", *cacheRoot)

	wg.Wait()
}
