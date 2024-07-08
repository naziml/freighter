package sqlite

import (
	"context"
	"path/filepath"

	"github.com/johnewart/freighter/server/layers"
	"zombiezen.com/go/log"
)

type RepositoryLayerStore struct {
	layers.RepositoryLayerStore
	store *DiskLayerFileStore
	DB    *DB
}

func NewRepositoryLayerStore(layerStore *DiskLayerFileStore, db *DB) *RepositoryLayerStore {
	return &RepositoryLayerStore{
		store: layerStore,
		DB:    db,
	}
}

func (s *RepositoryLayerStore) ListFiles(repo string, target string, path string) ([]FileRecord, error) {
	if files, err := s.DB.GetFilesForRepo(repo, target, path); err != nil {
		return nil, err
	} else {
		dirListing := make([]string, 0)
		for _, f := range files {
			dir, _ := filepath.Split(f.FilePath)
			dirListing = append(dirListing, dir)
		}

		result := make([]FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
		}

		for _, d := range dirListing {
			result = append(result, FileRecord{Name: d, Size: 0, IsDir: true})
		}
		return result, nil
	}
}

func (s *RepositoryLayerStore) ReadFile(repository string, target string, filename string) ([]byte, error) {
	ctx := context.Background()
	if f, err := s.DB.GetFileLayer(repository, target, filename); err != nil {
		return nil, err
	} else {

		log.Infof(ctx, "Fetching file %s from %s:%s in layer %s", filename, repository, target, f.LayerDigest)
		return s.store.ReadFile(f.LayerDigest, filename)
	}
}

func (s *RepositoryLayerStore) GetDirectoryTree(repository string, target string) []string {
	return s.DB.GetDirectoryTreeForRepo(repository, target)
}
