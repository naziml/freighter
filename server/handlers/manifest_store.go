package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/go-containerregistry/pkg/registry"
	"github.com/johnewart/freighter/server/storage"
	"zombiezen.com/go/log"
)

type ManifestConfig struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type ManifestLayer struct {
	MediaType string `json:"mediaType"`
	Size      int64  `json:"size"`
	Digest    string `json:"digest"`
}

type ContainerManifest struct {
	SchemaVersion int             `json:"schemaVersion"`
	MediaType     string          `json:"mediaType"`
	Config        ManifestConfig  `json:"config"`
	Layers        []ManifestLayer `json:"layers"`
}

type FreighterManifestStore struct {
	registry.ManifestStore
	lock  sync.RWMutex
	ctx   context.Context
	store storage.FreighterDataStore
}

func NewFreighterManifestStore(store storage.FreighterDataStore) *FreighterManifestStore {
	ctx := context.Background()

	return &FreighterManifestStore{
		ctx:   ctx,
		store: store,
	}
}

func (m *FreighterManifestStore) Get(repo string, target string) (*registry.Manifest, error) {

	if manifest, err := m.store.GetManifest(repo, target); err != nil {
		return nil, fmt.Errorf("manifest not found")
	} else {
		return &registry.Manifest{
			Repository: manifest.Repository,
			Target:     manifest.Target,
			MediaType:  manifest.MediaType,
			Blob:       manifest.RawBlob,
		}, nil
	}

}

func (m *FreighterManifestStore) Put(manifest registry.Manifest) error {

	log.Infof(m.ctx, "Put manifest %s:%s", manifest.Repository, manifest.Target)

	var cm ContainerManifest

	if err := json.Unmarshal(manifest.Blob, &cm); err != nil {
		return fmt.Errorf("failed to deserialize manifest blob: %v", err)
	}

	ms := storage.Manifest{
		Repository:      manifest.Repository,
		Target:          manifest.Target,
		MediaType:       manifest.MediaType,
		SchemaVersion:   cm.SchemaVersion,
		ConfigMediaType: cm.Config.MediaType,
		ConfigSize:      cm.Config.Size,
		ConfigDigest:    cm.Config.Digest,
		RawBlob:         manifest.Blob,
	}

	if mf, err := m.store.PutManifest(ms); err != nil {
		return fmt.Errorf("error creating manifest: %v", err)
	} else {

		for i, l := range cm.Layers {
			log.Infof(m.ctx, "Put layer %s:%s %s", manifest.Repository, manifest.Target, l.Digest)
			if err := m.store.PutLayer(storage.Layer{
				ManifestID: mf.ID,
				MediaType:  l.MediaType,
				Digest:     l.Digest,
				Size:       l.Size,
				Level:      i,
				Repository: manifest.Repository,
				Target:     manifest.Target,
			}); err != nil {
				return fmt.Errorf("error creating layer: %v", err)
			}
		}

		return nil
	}
}

func (m *FreighterManifestStore) Delete(repo string, target string) error {
	return m.store.DeleteManifest(repo, target)
}

func (m *FreighterManifestStore) GetTags(repo string) ([]string, error) {

	if repoManifests, err := m.store.ManifestsForRepo(repo); err != nil {
		return nil, fmt.Errorf("error finding manifests: %v", err)
	} else {
		tags := make([]string, 0, len(repoManifests))
		for _, manifest := range repoManifests {
			tags = append(tags, manifest.Target)
		}

		return tags, nil
	}
}

func (m *FreighterManifestStore) Exists(repo string, target string) bool {
	return m.store.ManifestExists(repo, target)
}

func (m *FreighterManifestStore) ListRepositories() []string {
	return m.store.ListRepositories()
}

func (m *FreighterManifestStore) ManifestsForRepository(repo string) ([]registry.Manifest, bool) {
	if repoManifests, err := m.store.ManifestsForRepo(repo); err != nil {
		return nil, false
	} else {
		result := make([]registry.Manifest, 0)
		for _, mf := range repoManifests {
			result = append(result, registry.Manifest{
				Repository: mf.Repository,
				Target:     mf.Target,
				MediaType:  mf.MediaType,
				Blob:       mf.RawBlob,
			})
		}
		return result, true
	}
}

var _ = (registry.ManifestStore)((*FreighterManifestStore)(nil))
