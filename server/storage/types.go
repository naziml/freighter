package storage

import (
	"context"
	"io"
	"strings"
)

type LayerFile struct {
	LayerDigest string
	FilePath    string
	Size        int64
	IsDir       bool
	Directory   string
}

func (lf LayerFile) Digest() Digest {
	return ParseDigest(lf.LayerDigest)
}

type Layer struct {
	ManifestID uint
	MediaType  string
	Size       int64
	Digest     string
	Level      int
	Repository string
	Target     string
}

type Manifest struct {
	ID              uint   `gorm:"primaryKey"`
	Repository      string `gorm:"index:idx_repotarget,unique"`
	Target          string `gorm:"index:idx_repotarget,unique"`
	SchemaVersion   int
	MediaType       string
	ConfigMediaType string
	ConfigDigest    string
	ConfigSize      int
	RawBlob         []byte
}

func (l Layer) GetDigest() Digest {
	return ParseDigest(l.Digest)
}

type Digest struct {
	Algorithm string
	Hash      string
}

func (d Digest) String() string {
	return d.Algorithm + ":" + d.Hash
}

func ParseDigest(s string) Digest {
	parts := strings.SplitN(s, ":", 2)
	return Digest{
		Algorithm: parts[0],
		Hash:      parts[1],
	}
}

type FreighterDataStore interface {
	ManifestsForRepo(repo string) ([]Manifest, error)
	GetManifest(repo string, target string) (*Manifest, error)
	PutManifest(Manifest) (Manifest, error)
	DeleteManifest(repo string, target string) error
	ManifestExists(repo string, target string) bool
	ListRepositories() []string
	PutLayer(Layer) error

	ListFiles(repo string, target string, path string) ([]FileRecord, error)
	GetDirectoryTreeForRepo(repository string, target string) []string
	ReadFile(repository string, target string, filename string) ([]byte, error)

	GetLayerReader(Digest) (io.ReadCloser, error)
	GetLayer(Digest) (*Layer, error)
	DeleteLayer(Digest) error

	IngestLayer(context.Context, Digest) error
	StoreBlob(digest Digest, rc io.ReadCloser) error

	String() string
}

type FileRecord struct {
	Name  string
	Size  int64
	IsDir bool
}

type LayerStore interface {
	StoreLayer(layerDigest string, data []byte) error
	LayerExists(repository, target, layerDigest string) bool
	ReadLayer(repository, target, layerDigest string) ([]byte, error)
	ReadFile(layerDigest string, filename string) ([]byte, error)
	IngestFiles(layerDigest string) error
}

type RepositoryStore interface {
	ListFiles(repo string, target string, path string) ([]FileRecord, error)
	ReadFile(repository string, target string, filename string) ([]byte, error)
	GetDirectoryTree(repository string, target string) []string
	DeleteLayer(layerDigest string) error
	IngestLayerFromReader(ctx context.Context, repository string, digestAlgorithm string, layerDigest string, r io.ReadCloser) error
}
