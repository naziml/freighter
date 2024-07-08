package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/johnewart/freighter/server/layers"
	"github.com/johnewart/freighter/server/layers/fs"

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
	LayerRepository layers.RepositoryStore
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

func (s *server) GetTree(ctx context.Context, in *pb.TreeRequest) (*pb.TreeReply, error) {
	log.Infof(ctx, "Received tree request for %s:%s", in.GetRepository(), in.GetTarget())

	files := make([]*pb.FileInfo, 0)

	if directories := s.LayerRepository.GetDirectoryTree(in.GetRepository(), in.GetTarget()); directories == nil {
		log.Errorf(ctx, "Error reading directory: %v", directories)
	} else {
		for _, d := range directories {
			files = append(files, &pb.FileInfo{Name: d, Size: 0, IsDir: true})
		}
	}

	return &pb.TreeReply{Files: files}, nil
}

func (s *server) GetDir(ctx context.Context, in *pb.DirRequest) (*pb.DirReply, error) {

	path := in.GetPath()
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	log.Infof(ctx, "Received dir request for %s:%s %v", in.GetRepository(), in.GetTarget(), path)

	files := make([]*pb.FileInfo, 0)

	if fileRecords, err := s.LayerRepository.ListFiles(in.GetRepository(), in.GetTarget(), path); err != nil {
		log.Errorf(ctx, "Error reading directory: %v", err)
	} else {
		log.Infof(ctx, "Found %d files in %s", len(fileRecords), path)
		for _, f := range fileRecords {
			log.Infof(ctx, "File: %s isdir: %s", f.Name, f.IsDir)
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
	db := fs.NewDB("manifests.db")

	if manifestStore, err := NewFreighterManifestStore(db); err != nil {
		log.Errorf(ctx, "Error creating manifest store: %v", err)
		return nil
	} else {
		log.Infof(ctx, "Created manifest store: %v", manifestStore)
		layerStore := fs.NewDiskLayerStore(rootPath, db)
		layerRepository := fs.NewDiskRepositoryStore(layerStore, db)
		indexingBlobstore := NewIndexingBlobStore(rootPath, layerRepository)

		registryHandler := registry.New(
			registry.WithWarning(.01, "Congratulations! You've won a lifetime's supply of free image pulls from this in-memory registry!"),
			registry.WithBlobHandler(indexingBlobstore),
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
