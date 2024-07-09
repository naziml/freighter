package fs

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnewart/freighter/server/storage"
	"zombiezen.com/go/log"
)

type DiskDataStore struct {
	storage.FreighterDataStore
	root string
	db   *DB
}

func NewDiskDataStore(rootPath string) (*DiskDataStore, error) {
	if db, err := NewDB(filepath.Join(rootPath, "metadata.db")); err != nil {
		return nil, err
	} else {

		return &DiskDataStore{
			root: rootPath,
			db:   db,
		}, nil
	}
}

func (d *DiskDataStore) GetManifest(repo string, target string) (*storage.Manifest, error) {
	return d.db.GetManifest(repo, target)
}

func (s *DiskDataStore) readFilesFromArchive(digest storage.Digest) ([]storage.FileRecord, error) {
	digestPath := s.getLayerPath(digest)
	log.Infof(context.Background(), "Reading files from layer: %s @ %s", digest, digestPath)
	files := make([]storage.FileRecord, 0)
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
			files = append(files, storage.FileRecord{Name: hdr.Name, Size: hdr.Size, IsDir: false})
		}
	}

	return files, nil
}

func (s *DiskDataStore) LayerFileNames(digest storage.Digest) ([]string, error) {
	files, err := s.readFilesFromArchive(digest)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(files))
	for _, f := range files {
		names = append(names, f.Name)
	}
	return names, nil
}

func (s *DiskDataStore) IngestLayer(ctx context.Context, digest storage.Digest) error {
	if files, err := s.readFilesFromArchive(digest); err != nil {
		log.Errorf(context.Background(), "Error reading files from layer: %v", err)
		return err
	} else {
		for _, r := range files {
			if !strings.HasPrefix(r.Name, "/") {
				r.Name = fmt.Sprintf("/%s", r.Name)
			}

			parts := strings.Split(r.Name, "/")
			dir := strings.Join(parts[:len(parts)-1], "/")

			if strings.HasSuffix(r.Name, "/") {
				r.IsDir = true
				if len(parts) > 2 {
					dir = strings.Join(parts[:len(parts)-2], "/")
				}
			}

			if dir == "" {
				dir = "/"
			}

			log.Infof(context.Background(), "Ingesting file: %s in '%s' (%v)", r.Name, dir, r.IsDir)

			if err := s.db.PutLayerFile(&storage.LayerFile{
				LayerDigest: digest.String(),
				FilePath:    r.Name,
				Size:        r.Size,
				IsDir:       r.IsDir,
				Directory:   dir,
			}); err != nil {
				log.Errorf(context.Background(), "Error creating layer file: %v", err)
				return err
			}
		}
		return nil
	}
}

func (s *DiskDataStore) readFile(digest storage.Digest, filename string) ([]byte, error) {
	ctx := context.Background()
	log.Infof(ctx, "Fetching file %s from  %s", filename, digest)

	filename = strings.TrimPrefix(filename, "/")
	tarFile := s.getLayerPath(digest)
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

func (s *DiskDataStore) GetDirectoryTreeForLayer(digest storage.Digest) ([]storage.FileRecord, error) {
	return s.readFilesFromArchive(digest)
}

func (s *DiskDataStore) DeleteLayer(digest storage.Digest) error {

	if err := s.db.DeleteLayer(digest); err != nil {
		log.Errorf(context.Background(), "Error deleting layer: %v", err)
		return err
	}

	return nil
}

func (s *DiskDataStore) getLayerPath(digest storage.Digest) string {
	return filepath.Join(s.root, digest.Algorithm, digest.Hash)
}

func (d *DiskDataStore) PutManifest(m storage.Manifest) (storage.Manifest, error) {
	return d.db.PutManifest(m)
}

func (d *DiskDataStore) GetLayer(digest storage.Digest) (*storage.Layer, error) {
	if layer, err := d.db.GetLayer(digest); err != nil {
		return nil, err
	} else {
		if layer == nil {
			layerPath := d.getLayerPath(digest)
			log.Infof(context.Background(), "Layer not found in DB: %s - checking if exists on disk at: %s", digest, layerPath)
			if stat, err := os.Stat(layerPath); err != nil {
				// No error, just doesn't exist
				return nil, nil
			} else {
				layer = &storage.Layer{
					Digest: digest.String(),
					Size:   stat.Size(),
				}
			}
		}

		return layer, nil
	}
}

func (s *DiskDataStore) ListFiles(repo string, target string, path string) ([]storage.FileRecord, error) {
	if files, err := s.db.GetFilesForRepo(repo, target, path); err != nil {
		return nil, err
	} else {
		result := make([]storage.FileRecord, 0, len(files))
		for _, f := range files {
			result = append(result, storage.FileRecord{Name: f.FilePath, Size: f.Size, IsDir: f.IsDir})
		}
		return result, nil
	}
}

func (s *DiskDataStore) ReadFile(repository string, target string, filename string) ([]byte, error) {
	ctx := context.Background()
	if lf, err := s.db.GetFileLayer(repository, target, filename); err != nil {
		return nil, err
	} else {
		log.Infof(ctx, "Fetching file %s from %s:%s in layer %s", filename, repository, target, lf.Digest())
		return s.readFile(lf.Digest(), filename)
	}
}

func (s *DiskDataStore) blobHashPath(algorithm string, h string) string {
	return filepath.Join(s.root, algorithm, h)
}

func (s *DiskDataStore) PutLayer(l storage.Layer) error {
	return s.db.PutLayer(&l)
}

func (s *DiskDataStore) StoreBlob(digest storage.Digest, rc io.ReadCloser) error {
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
	if err := os.MkdirAll(filepath.Join(s.root, digest.Algorithm), os.ModePerm); err != nil {
		return err
	}

	layerPath := s.getLayerPath(digest)
	if _, err := os.Stat(layerPath); err == nil {
		log.Infof(context.TODO(), "Layer already exists: %s", digest)
		os.Remove(layerPath)
	}

	err = os.Rename(f.Name(), layerPath)
	if err != nil {
		log.Errorf(context.TODO(), "Error renaming file: %v", err)
		return err
	}

	return nil
}

func (s *DiskDataStore) String() string {
	return fmt.Sprintf("DiskDataStore: %s", s.root)
}

var _ = (storage.FreighterDataStore)((*DiskDataStore)(nil))
