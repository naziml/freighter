package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/johnewart/freighter/server/handlers"
	"github.com/johnewart/freighter/server/storage"
	"github.com/johnewart/freighter/server/storage/types"

	"net"
	"os"

	pb "github.com/johnewart/freighter/freighter/proto"
	"google.golang.org/grpc"
	"zombiezen.com/go/log"

	"github.com/google/go-containerregistry/pkg/registry"
)

type server struct {
	pb.UnimplementedFreighterServer
	ManifestStore registry.ManifestStore
	DataStore     storage.FreighterDataStore
}

func (s *server) GetFile(ctx context.Context, in *pb.FileRequest) (*pb.FileReply, error) {
	filePath := in.GetPath()
	if filePath[:2] == "//" {
		filePath = filePath[1:]
	}

	log.Infof(ctx, "Fetching file %s from %s:%s", filePath, in.GetRepository(), in.GetTarget())
	if data, err := s.DataStore.ReadFile(in.GetRepository(), in.GetTarget(), filePath); err != nil {
		log.Errorf(ctx, "Error reading file: %v", err)
		return nil, err
	} else {
		return &pb.FileReply{Data: data}, nil
	}
}

func (s *server) GetTree(ctx context.Context, in *pb.TreeRequest) (*pb.TreeReply, error) {
	log.Infof(ctx, "Received tree request for %s:%s", in.GetRepository(), in.GetTarget())

	result := make([]*pb.FileInfo, 0)

	if files, err := s.DataStore.GetFilesForRepo(in.GetRepository(), in.GetTarget()); err != nil {
		log.Errorf(ctx, "Error reading tree: %v", err)
	} else {
		for _, f := range files {
			result = append(result, LayerFileToFileInfo(f))
		}
	}

	return &pb.TreeReply{Files: result}, nil
}

func (s *server) GetDir(ctx context.Context, in *pb.DirRequest) (*pb.DirReply, error) {

	path := in.GetPath()
	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	log.Infof(ctx, "Received dir request for %s:%s %v", in.GetRepository(), in.GetTarget(), path)

	files := make([]*pb.FileInfo, 0)

	if fileRecords, err := s.DataStore.ListFiles(in.GetRepository(), in.GetTarget(), path); err != nil {
		log.Errorf(ctx, "Error reading directory: %v", err)
	} else {
		log.Infof(ctx, "Found %d files in %s", len(fileRecords), path)
		for _, f := range fileRecords {
			log.Infof(ctx, "File: %s isdir: %v", f.FilePath, f.IsDir)
			files = append(files, LayerFileToFileInfo(f))
		}
	}

	return &pb.DirReply{Files: files}, nil
}

type FreighterServer struct {
	server          *grpc.Server
	registryHandler http.Handler
	ctx             context.Context
}

func NewFreighterServer(dataStore storage.FreighterDataStore) *FreighterServer {
	ctx := context.Background()
	s := grpc.NewServer()

	manifestStore := handlers.NewFreighterManifestStore(dataStore)
	indexingBlobstore := handlers.NewIndexingBlobStore(dataStore)

	registryHandler := registry.New(
		registry.WithWarning(.01, "Congratulations! You've won a lifetime's supply of free image pulls from this in-memory registry!"),
		registry.WithBlobHandler(indexingBlobstore),
		registry.WithManifestStore(manifestStore),
	)
	pb.RegisterFreighterServer(s, &server{
		DataStore:     dataStore,
		ManifestStore: manifestStore,
	})

	log.Infof(ctx, "Registering Freighter server with storage %s", dataStore.String())

	return &FreighterServer{
		ctx:             ctx,
		server:          s,
		registryHandler: registryHandler,
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

func LayerFileToFileInfo(lf types.LayerFile) *pb.FileInfo {
	fileType := pb.FileType_FILE
	extraData := ""
	switch lf.Type {
	case "F":
		fileType = pb.FileType_FILE
	case "D":
		fileType = pb.FileType_DIR
	case "S":
		fileType = pb.FileType_SYMLINK
		extraData = lf.ExtraData
	}

	return &pb.FileInfo{
		Name:       lf.FilePath,
		Path:       lf.Directory,
		Size:       uint64(lf.Size),
		IsDir:      lf.IsDir,
		Mode:       lf.Mode,
		ModTime:    uint64(lf.Mtime),
		AccessTime: uint64(lf.Atime),
		ChangeTime: uint64(lf.Ctime),
		Type:       fileType,
		ExtraData:  extraData,
	}
}
