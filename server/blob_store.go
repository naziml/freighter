package server

import (
	"context"
	"errors"
	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/johnewart/freighter/server/layers"
	"io"
	"os"
	"path/filepath"
	"zombiezen.com/go/log"
)

type IndexingBlobStore struct {
	registry.BlobHandler
	root string
	repo *layers.DiskLayerFileStore
}

func NewIndexingBlobStore(rootPath string, repo *layers.DiskLayerFileStore) *IndexingBlobStore {
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
	// Put the temp file in the same directory to avoid cross-device problems
	// during the os.Rename.  The filenames cannot conflict.
	f, err := os.CreateTemp(s.root, "upload-*")
	if err != nil {
		return err
	}

	if err := func() error {
		defer f.Close()
		_, err := io.Copy(f, rc)
		return err
	}(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(s.root, h.Algorithm), os.ModePerm); err != nil {
		return err
	}
	err = os.Rename(f.Name(), s.blobHashPath(h))
	if err != nil {
		log.Errorf(ctx, "Error renaming file: %v", err)
		return err
	}

	log.Infof(ctx, "Ingesting layer: %s", h.Hex)
	if err := s.repo.IngestFiles(h.Hex); err != nil {
		log.Errorf(ctx, "Error ingesting files: %v", err)
	}

	return nil
}

func (s *IndexingBlobStore) Delete(ctx context.Context, _ string, h v1.Hash) error {
	log.Infof(ctx, "Deleting layer: %s", h.Hex)
	if err := s.repo.DeleteLayer(h.Hex); err != nil {
		log.Errorf(ctx, "Error deleting layer: %v", err)
	}

	return os.Remove(s.blobHashPath(h))
}
