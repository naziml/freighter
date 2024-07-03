package layers

import (
	"context"

	"github.com/johnewart/freighter/server/data"
	"zombiezen.com/go/log"
)

type RepositoryLayerStore struct {
	store *DiskLayerFileStore
	DB    *data.DB
}

func NewRepositoryLayerStore(layerStore *DiskLayerFileStore, db *data.DB) *RepositoryLayerStore {
	return &RepositoryLayerStore{
		store: layerStore,
		DB:    db,
	}
}

func (s *RepositoryLayerStore) ListFiles(repo string, target string, path string) ([]FileRecord, error) {
	if files, err := s.DB.GetFilesForRepo(repo, target, path); err != nil {
		return nil, err
	} else {
		result := make([]FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
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
