package server

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/johnewart/freighter/server/layers"
	"zombiezen.com/go/log"
)

type IndexingBlobStore struct {
	registry.BlobHandler
	root string
	repo layers.RepositoryStore
}

func NewIndexingBlobStore(rootPath string, repo layers.RepositoryStore) *IndexingBlobStore {
	return &IndexingBlobStore{
		root: rootPath,
		repo: repo,
	}
}

func (s *IndexingBlobStore) blobHashPath(h v1.Hash) string {
	return filepath.Join(s.root, h.Algorithm, h.Hex)
}

func (s *IndexingBlobStore) Stat(_ context.Context, _ string, h v1.Hash) (int64, error) {
	fi, err := os.Stat(s.blobHashPath(h))
	if errors.Is(err, os.ErrNotExist) {
		return 0, registry.ErrBlobNotFound
	} else if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (s *IndexingBlobStore) Get(_ context.Context, _ string, h v1.Hash) (io.ReadCloser, error) {
	return os.Open(s.blobHashPath(h))
}

func (s *IndexingBlobStore) Put(ctx context.Context, repo string, h v1.Hash, rc io.ReadCloser) error {

	s.repo.IngestLayerFromReader(ctx, repo, h.Algorithm, h.Hex, rc)
	return nil
}

func (s *IndexingBlobStore) Delete(ctx context.Context, _ string, h v1.Hash) error {
	log.Infof(ctx, "Deleting layer: %s", h.Hex)
	if err := s.repo.DeleteLayer(h.Hex); err != nil {
		log.Errorf(ctx, "Error deleting layer: %v", err)
	}

	return os.Remove(s.blobHashPath(h))
}
