package server

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"zombiezen.com/go/log"
)

type LayerRepository struct {
	registry.BlobHandler
	RootPath string
	DB       *DB
}

type FileRecord struct {
	Name  string
	Size  int64
	IsDir bool
}

func NewLayerRepository(rootPath string, db *DB) *LayerRepository {
	return &LayerRepository{
		RootPath: rootPath,
		DB:       db,
	}
}

func (l *LayerRepository) GetLayerPath(digest string) string {
	return filepath.Join(l.RootPath, "sha256", digest)
}

func (l *LayerRepository) readFilesFromArchive(digest string) ([]FileRecord, error) {
	digestPath := l.GetLayerPath(digest)
	log.Infof(context.Background(), "Reading files from layer: %s @ %s", digest, digestPath)
	files := make([]FileRecord, 0)
	if file, err := os.Open(digestPath); err == nil {
		archive, err := gzip.NewReader(file)

		if err != nil {
			return nil, err
		}

		tr := tar.NewReader(archive)

		if tr == nil {
			return nil, err
		}

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			files = append(files, FileRecord{Name: hdr.Name, Size: hdr.Size, IsDir: false})
		}
	}

	return files, nil
}

func (l *LayerRepository) IngestFiles(digest string) error {
	if files, err := l.readFilesFromArchive(digest); err != nil {
		log.Errorf(context.Background(), "Error reading files from layer: %v", err)
		return err
	} else {
		for _, r := range files {
			log.Infof(context.Background(), "Ingesting file: %s", r.Name)
			if err := l.DB.PutLayerFile(&LayerFile{
				LayerDigest: digest,
				FilePath:    r.Name,
				Size:        r.Size,
				IsDir:       r.IsDir,
			}); err != nil {
				log.Errorf(context.Background(), "Error creating layer file: %v", err)
				return err
			}
		}
		return nil
	}
}

func (l *LayerRepository) ListFiles(repo string, target string) ([]FileRecord, error) {
	if files, err := l.DB.GetFilesForRepo(repo, target); err != nil {
		return nil, err
	} else {
		result := make([]FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
		}
		return result, nil
	}
}

func (l *LayerRepository) ReadFile(repository string, target string, filename string) ([]byte, error) {
	ctx := context.Background()
	if f, err := l.DB.GetFileLayer(repository, target, filename); err != nil {
		return nil, err
	} else {
		log.Infof(ctx, "Fetching file %s from %s:%s in layer %s", filename, repository, target, f.LayerDigest)

		tarFile := l.GetLayerPath(f.LayerDigest)
		file, err := os.Open(tarFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		archive, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer archive.Close()

		tr := tar.NewReader(archive)

		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if hdr.Name == filename {
				buf, err := ioutil.ReadAll(tr)
				if err != nil {
					return nil, err
				}
				return buf, nil
			}
		}
		return nil, fmt.Errorf("File not found: %s", filename)
	}
}

func (l *LayerRepository) blobHashPath(h v1.Hash) string {
	return filepath.Join(l.RootPath, h.Algorithm, h.Hex)
}

func (l *LayerRepository) Stat(_ context.Context, _ string, h v1.Hash) (int64, error) {
	fi, err := os.Stat(l.blobHashPath(h))
	if errors.Is(err, os.ErrNotExist) {
		return 0, registry.ErrBlobNotFound
	} else if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func (l *LayerRepository) Get(_ context.Context, _ string, h v1.Hash) (io.ReadCloser, error) {
	return os.Open(l.blobHashPath(h))
}

func (l *LayerRepository) Put(ctx context.Context, repo string, h v1.Hash, rc io.ReadCloser) error {
	// Put the temp file in the same directory to avoid cross-device problems
	// during the os.Rename.  The filenames cannot conflict.
	f, err := os.CreateTemp(l.RootPath, "upload-*")
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
	if err := os.MkdirAll(filepath.Join(l.RootPath, h.Algorithm), os.ModePerm); err != nil {
		return err
	}
	err = os.Rename(f.Name(), l.blobHashPath(h))
	if err != nil {
		log.Errorf(ctx, "Error renaming file: %v", err)
		return err
	}

	log.Infof(ctx, "Ingesting layer: %s", h.Hex)
	if err := l.IngestFiles(h.Hex); err != nil {
		log.Errorf(ctx, "Error ingesting files: %v", err)
	}

	return nil
}

func (l *LayerRepository) Delete(ctx context.Context, _ string, h v1.Hash) error {
	log.Infof(ctx, "Deleting layer: %s", h.Hex)
	l.DB.DeleteLayer(h.Hex)
	return os.Remove(l.blobHashPath(h))
}
