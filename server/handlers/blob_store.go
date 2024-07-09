package handlers

import (
	"compress/gzip"
	"context"
	"io"

	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/johnewart/freighter/server/storage"
	"zombiezen.com/go/log"
)

type IndexingBlobStore struct {
	registry.BlobHandler
	store storage.FreighterDataStore
}

func NewIndexingBlobStore(store storage.FreighterDataStore) *IndexingBlobStore {
	return &IndexingBlobStore{
		store: store,
	}
}

func (s *IndexingBlobStore) Stat(_ context.Context, _ string, h v1.Hash) (int64, error) {
	layer, err := s.store.GetLayer(storage.Digest{Algorithm: h.Algorithm, Hash: h.Hex})
	if err != nil {
		log.Errorf(context.Background(), "Error stating layer: %v", err)
		return 0, err
	} else if layer == nil {
		log.Errorf(context.Background(), "Layer not found: %s", h.Hex)
		return 0, registry.ErrBlobNotFound
	}

	log.Infof(context.Background(), "Stat layer: %s:%s with size %d", h.Algorithm, h.Hex, layer.Size)
	return layer.Size, nil
}

func (s *IndexingBlobStore) Get(_ context.Context, _ string, h v1.Hash) (io.ReadCloser, error) {
	rc, err := s.store.GetLayerReader(storage.Digest{Algorithm: h.Algorithm, Hash: h.Hex})
	if err != nil {
		log.Errorf(context.Background(), "Error getting layer: %v", err)
		return nil, err
	}

	return rc, nil
}

func (s *IndexingBlobStore) Put(ctx context.Context, repo string, h v1.Hash, rc io.ReadCloser) error {
	log.Infof(ctx, "Storing layer data: %s:%s", h.Algorithm, h.Hex)
	if err := s.store.StoreBlob(storage.Digest{Algorithm: h.Algorithm, Hash: h.Hex}, rc); err != nil {
		log.Errorf(ctx, "Error storing file: %v", err)
		return err
	}

	if _, err := gzip.NewReader(rc); err != nil {
		log.Infof(ctx, "Layer does not look like a gzip archive, skipping ingestion: %s:%s", h.Algorithm, h.Hex)
		return nil
	} else {
		log.Infof(ctx, "Ingesting layer: %s:%s", h.Algorithm, h.Hex)
		err := s.store.IngestLayer(ctx, storage.Digest{Algorithm: h.Algorithm, Hash: h.Hex})
		if err != nil {
			log.Errorf(ctx, "Error ingesting layer: %v", err)
		} else {
			log.Infof(ctx, "Ingested layer: %s:%s", h.Algorithm, h.Hex)
		}

		return err
	}

}

func (s *IndexingBlobStore) Delete(ctx context.Context, _ string, h v1.Hash) error {
	log.Infof(ctx, "Deleting layer: %s:%s", h.Algorithm, h.Hex)

	if err := s.store.DeleteLayer(storage.Digest{Algorithm: h.Algorithm, Hash: h.Hex}); err != nil {
		log.Errorf(ctx, "Error deleting layer: %v", err)
		return err
	}

	return nil
}
