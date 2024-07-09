package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"sync"

	"github.com/johnewart/freighter/server"
	"github.com/johnewart/freighter/server/storage/fs"
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

	if _, err := os.Stat(*cacheRoot); errors.Is(err, os.ErrNotExist) {
		log.Infof(context.Background(), "Cache root %s does not exist, creating", *cacheRoot)
		os.MkdirAll(*cacheRoot, 0755)
	}

	ctx := context.Background()

	dataStore, err := fs.NewDiskDataStore(*cacheRoot)
	if err != nil {
		log.Errorf(ctx, "Error creating data store: %v", err)
		return
	}

	fs := server.NewFreighterServer(dataStore)
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
