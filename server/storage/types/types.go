package types

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

type FileRecord struct {
	Name  string
	Size  int64
	IsDir bool
}

type LayerStore interface {
	ReadFile(LayerFile) ([]byte, error)
	StoreLayerBlob(Digest, io.ReadCloser) error
	IngestLayer(context.Context, Digest) ([]LayerFile, error)
	DeleteLayer(Digest) error
	GetLayer(Digest) (*Layer, error)
	String() string
	GetLayerReader(Digest) (io.ReadCloser, error)
}

type MetadataStore interface {
	GetManifest(repository, target string) (*Manifest, error)
	DeleteLayer(layerDigest Digest) error
	StoreManifest(Manifest) (Manifest, error)
	DeleteManifest(repository, target string) error
	GetFilesForRepo(repository, target, path string) ([]LayerFile, error)
	GetDirectoryTreeForRepo(repository, target string) []string
	GetLayerFile(repository, target, filename string) (LayerFile, error)
	GetDirectoryTreeForLayer(digest Digest) ([]FileRecord, error)
	StoreLayerFiles([]LayerFile) error
	StoreLayer(layer Layer) error
	ManifestsForRepo(repository string) ([]Manifest, error)
	ListRepositories() []string
	String() string
}
