package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"net"
	"os"

	pb "github.com/johnewart/freighter/freighter/proto"
	"google.golang.org/grpc"
	"zombiezen.com/go/log"

	"github.com/google/go-containerregistry/pkg/registry"
)

type server struct {
	pb.UnimplementedFreighterServer
	ManifestStore   registry.ManifestStore
	LayerRepository *LayerRepository
}

func (s *server) GetFile(ctx context.Context, in *pb.FileRequest) (*pb.FileReply, error) {
	log.Infof(ctx, "Fetching file %s from %s:%s", in.GetPath(), in.GetRepository(), in.GetTarget())
	if data, err := s.LayerRepository.ReadFile(in.GetRepository(), in.GetTarget(), in.GetPath()); err != nil {
		log.Errorf(ctx, "Error reading file: %v", err)
		return nil, err
	} else {
		return &pb.FileReply{Data: data}, nil
	}
}

func (s *server) GetDir(ctx context.Context, in *pb.DirRequest) (*pb.DirReply, error) {
	log.Infof(ctx, "Received dir request for %s:%s /%v", in.GetRepository(), in.GetTarget(), in.GetPath())

	files := make([]*pb.FileInfo, 0)

	if fileRecords, err := s.LayerRepository.ListFiles(in.GetRepository(), in.GetTarget()); err != nil {
		log.Errorf(ctx, "Error reading directory: %v", err)
	} else {
		log.Infof(ctx, "Found %d files", len(fileRecords))
		for _, f := range fileRecords {
			files = append(files, &pb.FileInfo{Name: f.Name, Size: f.Size, IsDir: f.IsDir})
		}
	}

	return &pb.DirReply{Files: files}, nil
}

type FreighterServer struct {
	server          *grpc.Server
	registryHandler http.Handler
	ctx             context.Context
}

func NewFreighterServer(rootPath string) *FreighterServer {
	ctx := context.Background()
	s := grpc.NewServer()
	db := NewDB("manifests.db")

	if manifestStore, err := NewFreighterManifestStore(db); err != nil {
		log.Errorf(ctx, "Error creating manifest store: %v", err)
		return nil
	} else {
		log.Infof(ctx, "Created manifest store: %v", manifestStore)
		layerRepository := NewLayerRepository(rootPath, db)

		registryHandler := registry.New(
			registry.WithWarning(.01, "Congratulations! You've won a lifetime's supply of free image pulls from this in-memory registry!"),
			registry.WithBlobHandler(layerRepository),
			registry.WithManifestStore(manifestStore),
		)
		pb.RegisterFreighterServer(s, &server{
			LayerRepository: layerRepository,
			ManifestStore:   manifestStore,
		})

		log.Infof(ctx, "Registering Freighter server with layer root at %s", rootPath)

		return &FreighterServer{
			ctx:             ctx,
			server:          s,
			registryHandler: registryHandler,
		}
	}
}

func (fs *FreighterServer) ServeRegistry(host string, port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	porti := listener.Addr().(*net.TCPAddr).Port
	log.Infof(fs.ctx, "serving on port %d", porti)

	s := &http.Server{
		ReadHeaderTimeout: 5 * time.Second, // prevent slowloris, quiet linter
		Handler:           fs.registryHandler,
	}
	s.Serve(listener)
	return nil
}

func (fs *FreighterServer) Serve(host string, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Errorf(fs.ctx, "failed to listen: %v", err)
		os.Exit(1)
	}
	log.Infof(fs.ctx, "server listening at %v", lis.Addr())
	if err := fs.server.Serve(lis); err != nil {
		log.Errorf(fs.ctx, "failed to serve: %v", err)
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
