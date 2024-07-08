package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/johnewart/freighter/server/layers"
	"zombiezen.com/go/log"
)

type DiskRepositoryStore struct {
	root  string
	store *DiskLayerStore
	DB    *DB
}

func NewDiskRepositoryStore(layerStore *DiskLayerStore, db *DB) *DiskRepositoryStore {
	return &DiskRepositoryStore{
		root:  layerStore.RootPath,
		store: layerStore,
		DB:    db,
	}
}

func (s *DiskRepositoryStore) ListFiles(repo string, target string, path string) ([]layers.FileRecord, error) {
	if files, err := s.DB.GetFilesForRepo(repo, target, path); err != nil {
		return nil, err
	} else {
		result := make([]layers.FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, layers.FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
		}
		return result, nil
	}
}

func (s *DiskRepositoryStore) ReadFile(repository string, target string, filename string) ([]byte, error) {
	ctx := context.Background()
	if f, err := s.DB.GetFileLayer(repository, target, filename); err != nil {
		return nil, err
	} else {

		log.Infof(ctx, "Fetching file %s from %s:%s in layer %s", filename, repository, target, f.LayerDigest)
		return s.store.ReadFile(f.LayerDigest, filename)
	}
}

func (s *DiskRepositoryStore) GetDirectoryTree(repository string, target string) []string {
	return s.DB.GetDirectoryTreeForRepo(repository, target)
}

func (s *DiskRepositoryStore) DeleteLayer(layerDigest string) error {
	return s.DB.DeleteLayer(layerDigest)
}

func (s *DiskRepositoryStore) blobHashPath(algorithm string, h string) string {
	return filepath.Join(s.root, algorithm, h)
}

func (s *DiskRepositoryStore) IngestLayerFromReader(ctx context.Context, repository string, digestAlgorithm string, layerDigest string, rc io.ReadCloser) error {
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
	if err := os.MkdirAll(filepath.Join(s.root, digestAlgorithm), os.ModePerm); err != nil {
		return err
	}
	err = os.Rename(f.Name(), s.blobHashPath(digestAlgorithm, layerDigest))
	if err != nil {
		log.Errorf(ctx, "Error renaming file: %v", err)
		return err
	}

	log.Infof(ctx, "Ingesting layer: %s", layerDigest)

	if err := s.store.IngestFilesForLayer(layerDigest); err != nil {
		log.Errorf(ctx, "Error ingesting files: %v", err)
		return err
	}

	return nil
}
