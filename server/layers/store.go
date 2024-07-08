package layers

import (
	"context"
	"io"
)

type FileRecord struct {
	Name  string
	Size  int64
	IsDir bool
}

type LayerStore interface {
	StoreLayer(layerDigest string, data []byte) error
	LayerExists(repository, target, layerDigest string) bool
	ReadLayer(repository, target, layerDigest string) ([]byte, error)
	IngestFiles(layerDigest string) error
}

type RepositoryStore interface {
	ListFiles(repo string, target string, path string) ([]FileRecord, error)
	ReadFile(repository string, target string, filename string) ([]byte, error)
	GetDirectoryTree(repository string, target string) []string
	DeleteLayer(layerDigest string) error
	IngestLayerFromReader(ctx context.Context, repository string, digestAlgorithm string, layerDigest string, r io.ReadCloser) error
}
