/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"

	pb "github.com/johnewart/freighter/freighter/proto"
	ffs "github.com/johnewart/freighter/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"zombiezen.com/go/log"
)

const (
	defaultName = "world"
)

var (
	addr   = flag.String("addr", "localhost:50051", "the address to connect to")
	debug  = flag.Bool("debug", false, "print debug data")
	repo   = flag.String("repo", "", "container repository")
	target = flag.String("target", "", "the container target")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf(nil, "did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFreighterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Errorf(ctx, "Usage: hello MOUNTPOINT")
	}
	opts := &fs.Options{}
	opts.Debug = *debug

	log.Infof(ctx, "Mounting %s:%s!", *repo, *target)
	//root := ffs.NewFreighterTree(ctx, c, *repo, *target, "")

	root := &ffs.FreighterRoot{
		Client:     c,
		Repository: *repo,
		Target:     *target,
	}
	server, err := fs.Mount(flag.Arg(0), root, opts)
	if err != nil {
		log.Errorf(ctx, "Mount fail: %v", err)
	}
	//root.LoadTree(ctx)
	server.Wait()

}
