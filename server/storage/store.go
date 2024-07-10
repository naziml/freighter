package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/johnewart/freighter/server/storage/types"
	"zombiezen.com/go/log"
)

type FreighterDataStore struct {
	metadata types.MetadataStore
	layers   types.LayerStore
}

func NewFreighterDataStore(metadataStore types.MetadataStore, fileStore types.LayerStore) (FreighterDataStore, error) {
	return FreighterDataStore{
		metadata: metadataStore,
		layers:   fileStore,
	}, nil
}

func (s *FreighterDataStore) GetManifest(repo string, target string) (*types.Manifest, error) {
	return s.metadata.GetManifest(repo, target)
}

func (s *FreighterDataStore) IngestLayer(ctx context.Context, digest types.Digest) error {
	if layerFiles, err := s.layers.IngestLayer(ctx, digest); err != nil {
		log.Errorf(ctx, "Error ingesting layer: %v", err)
		return err
	} else {
		return s.metadata.StoreLayerFiles(layerFiles)
	}
}

func (s *FreighterDataStore) GetDirectoryTreeForLayer(digest types.Digest) ([]types.FileRecord, error) {
	return s.metadata.GetDirectoryTreeForLayer(digest)
}

func (s *FreighterDataStore) GetDirectoryTreeForRepo(repo string, target string) []string {
	return s.metadata.GetDirectoryTreeForRepo(repo, target)
}

func (s *FreighterDataStore) GetFilesForRepo(repo string, target string) ([]types.FileRecord, error) {
	if files, err := s.metadata.GetFilesForRepo(repo, target, ""); err != nil {
		return nil, err
	} else {
		result := make([]types.FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, types.FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
		}
		return result, nil
	}

	//return s.metadata.GetFilesForRepo(repo, target, "")

}

func (s *FreighterDataStore) DeleteLayer(digest types.Digest) error {

	if err := s.metadata.DeleteLayer(digest); err != nil {
		log.Errorf(context.Background(), "Error deleting layer: %v", err)
		return err
	}

	if err := s.layers.DeleteLayer(digest); err != nil {
		log.Errorf(context.Background(), "Error deleting layer: %v", err)
		return err
	}

	return nil
}

func (s *FreighterDataStore) PutManifest(m types.Manifest) (types.Manifest, error) {
	return s.metadata.StoreManifest(m)
}

func (s *FreighterDataStore) GetLayer(digest types.Digest) (*types.Layer, error) {
	if layer, err := s.layers.GetLayer(digest); err != nil {
		return nil, err
	} else {
		if layer != nil {
			return layer, nil
		} else {
			return s.layers.GetLayer(digest)
		}
	}
}

func (s *FreighterDataStore) ListFiles(repo string, target string, path string) ([]types.FileRecord, error) {
	if files, err := s.metadata.GetFilesForRepo(repo, target, path); err != nil {
		return nil, err
	} else {
		result := make([]types.FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, types.FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
		}
		return result, nil
	}
}

func (s *FreighterDataStore) ReadFile(repository string, target string, filename string) ([]byte, error) {
	ctx := context.Background()
	if lf, err := s.metadata.GetLayerFile(repository, target, filename); err != nil {
		return nil, err
	} else {
		log.Infof(ctx, "Fetching file %s from %s:%s in layer %s", filename, repository, target, lf.Digest())
		return s.layers.ReadFile(*lf)
	}
}

func (s *FreighterDataStore) PutLayer(l types.Layer) error {
	return s.metadata.StoreLayer(l)
}

func (s *FreighterDataStore) StoreLayerBlob(digest types.Digest, rc io.ReadCloser) error {
	return s.layers.StoreLayerBlob(digest, rc)
}
func (s *FreighterDataStore) String() string {
	return fmt.Sprintf("FreighterDataStore - metadata: %s layers: %s", s.metadata.String(), s.layers.String())
}

func (s *FreighterDataStore) ManifestsForRepo(repository string) ([]types.Manifest, error) {
	return s.metadata.ManifestsForRepo(repository)
}

func (s *FreighterDataStore) DeleteManifest(repo string, target string) error {
	return s.metadata.DeleteManifest(repo, target)
}

func (s *FreighterDataStore) ManifestExists(repo string, target string) bool {
	if m, err := s.metadata.GetManifest(repo, target); err != nil {
		return false
	} else {
		return m != nil
	}
}

func (s *FreighterDataStore) GetLayerReader(digest types.Digest) (io.ReadCloser, error) {
	return s.layers.GetLayerReader(digest)
}

func (s *FreighterDataStore) ListRepositories() []string {
	return s.metadata.ListRepositories()
}
